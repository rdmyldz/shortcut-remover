package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
typedef struct _CONSOLE_SCREEN_BUFFER_INFO {
  COORD      dwSize;
  COORD      dwCursorPosition;
  WORD       wAttributes;
  SMALL_RECT srWindow;
  COORD      dwMaximumWindowSize;
} CONSOLE_SCREEN_BUFFER_INFO;

type ConsoleScreenBufferInfo struct {
	Size              Coord
	CursorPosition    Coord
	Attributes        uint16
	Window            SmallRect
	MaximumWindowSize Coord
}
*/

/*
func GetConsoleScreenBufferInfo(console Handle, info *ConsoleScreenBufferInfo) (err error) {
	r1, _, e1 := syscall.Syscall(procGetConsoleScreenBufferInfo.Addr(), 2, uintptr(console), uintptr(unsafe.Pointer(info)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
*/

/*
BOOL WINAPI SetConsoleScreenBufferSize(
  _In_ HANDLE hConsoleOutput,
  _In_ COORD  dwSize
);

*/

// var modkernel32 = NewLazySystemDLL("kernel32.dll")

// const procGetConsoleScreenBufferInfo = modkernel32.NewProc("GetConsoleScreenBufferInfo")

var krnl32Mod = syscall.MustLoadDLL("kernel32.dll")
var procSetConsoleSize = krnl32Mod.MustFindProc("SetConsoleScreenBufferSize")

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWinsize() (*winsize, error) {
	ws := new(winsize)
	fd := os.Stdout.Fd()
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(fd), &info); err != nil {
		return nil, err
	}

	ws.Col = uint16(info.Window.Right - info.Window.Left + 1)
	ws.Row = uint16(info.Window.Bottom - info.Window.Top + 1)

	return ws, nil
}

var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procSetConsoleScreen = modkernel32.NewProc("GetConsoleScreenBufferInfo")

// Global screen buffer
// Its not recommended write to buffer dirrectly, use package Print,Printf,Println fucntions instead.
var Screen *bytes.Buffer = new(bytes.Buffer)

// Move cursor to given position
func MoveCursor(x int, y int) {
	fmt.Fprintf(Screen, "\033[%d;%dH", y, x)
}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(Screen, a...)
}

// Flush buffer and ensure that it will not overflow screen
func Flush() {
	for _, str := range strings.SplitAfter(Screen.String(), "\n") {
		// if idx > Height() {
		// 	return
		// }

		Output.WriteString(str)
	}

	Output.Flush()
	Screen.Reset()
}

var Output *bufio.Writer = bufio.NewWriter(os.Stdout)

// Clear screen
func Clear() {
	// Output.WriteString("\033[H\033[2J")
	Output.WriteString("\033[2J")
}

func getCharNumber(num int) string {
	var s strings.Builder
	for i := 0; i < 35; i++ {
		for i := 0; i < num; i++ {
			s.WriteString("a")
		}
	}

	return s.String()
}

type coord struct {
	dwSize windows.Coord
}

func SetConsoleScreenBufferInfo(console windows.Handle, info *coord) (err error) {
	r1, _, e1 := syscall.Syscall(procSetConsoleScreen.Addr(), 2, uintptr(console), uintptr(unsafe.Pointer(info)), 0)
	if r1 == 0 {
		return e1
	}
	return
}

var (
	kernel32Dll    *syscall.LazyDLL  = syscall.NewLazyDLL("Kernel32.dll")
	setConsoleMode *syscall.LazyProc = kernel32Dll.NewProc("SetConsoleMode")
)

func EnableVirtualTerminalProcessing(stream syscall.Handle, enable bool) error {
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING uint32 = 0x4

	var mode uint32
	err := syscall.GetConsoleMode(syscall.Stdout, &mode)
	if err != nil {
		return err
	}

	if enable {
		mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	} else {
		mode &^= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	}

	ret, _, err := setConsoleMode.Call(uintptr(stream), uintptr(mode))
	if ret == 0 {
		return err
	}

	return nil
}

func writeNums() <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	return ch
}

func main() {
	EnableVirtualTerminalProcessing(syscall.Stdout, true)
	log.Println("deneme")

	// dwSize := windows.Coord{X: 120, Y: 50}

	// s := getCharNumber(125)
	// fmt.Printf("s: %s\n", s)
	// w, h, err := term.GetSize(int(os.Stdout.Fd()))
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Printf("w: %d, h: %d\n", w, h)

	// wSize, err := getWinsize()
	// if err != nil {
	// 	log.Fatalf("in getWinsize(), %v\n", err)
	// }

	// log.Printf("wsize: %v\n", wSize)

	// // fd := os.Stdout.Fd()

	// r1, r2, err := procSetConsoleSize.Call(uintptr(os.Stdout.Fd()), uintptr(unsafe.Pointer(&dwSize)))
	// fmt.Printf("r1: %v --- r2: %v ---- err: %v\n", r1, r2, err)

	// scanner := bufio.NewScanner(os.Stdin)

	// scanner.Scan()

	ch := writeNums()

	Clear() // Clear current screen

	for num := range ch {

		// By moving cursor to top-left position we ensure that console output
		// will be overwritten each time, instead of adding new.
		MoveCursor(1, 1)

		// Println("Current Time: ", time.Now().Format(time.RFC1123))
		Println("another row")
		Println("Current Time: ", num)

		Flush() // Call it every time at the end of rendering

		// time.Sleep(time.Second)

	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

}
