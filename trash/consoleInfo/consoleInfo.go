package main

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// StringToCharPtr converts a Go string into pointer to a null-terminated cstring.
// This assumes the go string is already ANSI encoded.
func StringToCharPtr(str string) *uint8 {
	chars := append([]byte(str), 0) // null terminated
	return &chars[0]
}

// StringToUTF16Ptr converts a Go string into a pointer to a null-terminated UTF-16 wide string.
// This assumes str is of a UTF-8 compatible encoding so that it can be re-encoded as UTF-16.
func StringToUTF16Ptr(str string) *uint16 {
	wchars := utf16.Encode([]rune(str + "\x00"))
	return &wchars[0]
}

// func CreateJobObject(attr *syscall.SecurityAttributes, name string) (syscall.Handle, error) {
// 	r1, _, err := procCreateJobA.Call(
// 		uintptr(unsafe.Pointer(attr)),
// 		uintptr(unsafe.Pointer(StringToCharPtr(name))),
// 	)

// 	if err != syscall.Errno(0) {
// 		return 0, err
// 	}

// 	return syscall.Handle(r1), nil

// }

/*
BOOL SystemParametersInfoW(
  [in]      UINT  uiAction,
  [in]      UINT  uiParam,
  [in, out] PVOID pvParam,
  [in]      UINT  fWinIni
);
*/

var (
	user32DLL = windows.NewLazyDLL("User32.dll")
	procUser  = user32DLL.NewProc("SystemParametersInfoW")

	kernel32DLL = windows.NewLazyDLL("Kernel32.dll")
	// procKernel32 = kernel32DLL.NewProc("GetConsoleScreenBufferInfoEx")
	// procKernel32 = kernel32DLL.NewProc("SetConsoleScreenBufferInfoEx")
	procKernel32 = kernel32DLL.NewProc("SetConsoleWindowInfo")
	// winSpool     = windows.LazyDLL("Winspool.drv")
	winSpool          = syscall.MustLoadDLL("Winspool.drv")
	getDefaultPrinter = winSpool.MustFindProc("GetDefaultPrinterW")
)

/*
typedef struct _CONSOLE_SCREEN_BUFFER_INFO {
  COORD      dwSize;
  COORD      dwCursorPosition;
  WORD       wAttributes;
  SMALL_RECT srWindow;
  COORD      dwMaximumWindowSize;
} CONSOLE_SCREEN_BUFFER_INFO;
*/
type conBufInfo struct {
	dwSize              Coord
	dwCursorPosition    Coord
	wAttributes         uint16
	srWindow            SmallRect
	dwMaximumWindowSize Coord
}

/*
typedef struct _COORD {
  SHORT X;
  SHORT Y;
} COORD, *PCOORD;

typedef struct _SMALL_RECT {
  SHORT Left;
  SHORT Top;
  SHORT Right;
  SHORT Bottom;
} SMALL_RECT;

*/

type Coord struct {
	x int16
	y int16
}

type SmallRect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

type ColorTable uint32

/*
typedef struct _CONSOLE_SCREEN_BUFFER_INFOEX {
  ULONG      cbSize;
  COORD      dwSize;
  COORD      dwCursorPosition;
  WORD       wAttributes;
  SMALL_RECT srWindow;
  COORD      dwMaximumWindowSize;
  WORD       wPopupAttributes;
  BOOL       bFullscreenSupported;
  COLORREF   ColorTable[16];
} CONSOLE_SCREEN_BUFFER_INFOEX, *PCONSOLE_SCREEN_BUFFER_INFOEX;
*/

type conBufInfoEx struct {
	cbSize               uint32
	dwSize               Coord
	dwCursorPosition     Coord
	wAttributes          uint16
	srWindow             SmallRect
	dwMaximumWindowSize  Coord
	wPopupAttributes     uint16
	bFullscreenSupported int32
	colorTable           []ColorTable
}

/*
BOOL WINAPI SetConsoleWindowInfo(
  _In_       HANDLE     hConsoleOutput,
  _In_       BOOL       bAbsolute,
  _In_ const SMALL_RECT *lpConsoleWindow
);
*/

/*
typedef struct _SHARE_INFO_1 {
  LMSTR shi1_netname;
  DWORD shi1_type;
  LMSTR shi1_remark;
} SHARE_INFO_1, *PSHARE_INFO_1, *LPSHARE_INFO_1;
*/

type shareInfo struct {
	shi1_netname string
	shi1_type    uint32
	shi1_remark  string
}

func main() {
	// imagePath := `C:\Windows\Web\Wallpaper\Theme1\img4.jpg`
	// ip, err := windows.UTF16PtrFromString(imagePath)
	// fmt.Printf("err: %v\n", err)
	// r1, r2, err := procUser.Call(20, 0, uintptr(unsafe.Pointer(ip)), 0x001A)
	// fmt.Printf("r1: %v --- r2: %v ---- err: %v\n", r1, r2, err)

	// var deneme conBufInfoEx
	/*
		winsize := SmallRect{
			Left:   5,
			Top:    5,
			Right:  55,
			Bottom: 5,
		}

		stdhandle, err := syscall.GetStdHandle(4294967285)
		fmt.Printf("error: %v\n", err)

		// r1, r2, err := procKernel32.Call(uintptr(stdhandle), uintptr(unsafe.Pointer(&deneme)))
		// r1, r2, err := procKernel32.Call(uintptr(stdhandle), uintptr(unsafe.Pointer(&deneme)))
		r1, r2, err := procKernel32.Call(uintptr(stdhandle), uintptr(0), uintptr(unsafe.Pointer(&winsize)))
		fmt.Printf("r1: %v --- r2: %v ---- err: %v\n", r1, r2, err)

		fmt.Printf("deneme: %+v\n", winsize)
		log.Println("hey hey")
		fmt.Scanln()
	*/
	var bufSize uint32
	r1, r2, err := getDefaultPrinter.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(&bufSize)),
	)
	fmt.Printf("r1: %v --- r2: %v ---- err: %v\n", r1, r2, err)
	fmt.Printf("bufsize: %+v\n", bufSize)

	buf := make([]uint16, bufSize)
	r1, r2, err = getDefaultPrinter.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufSize)),
	)
	fmt.Printf("r1: %v --- r2: %v ---- err: %v\n", r1, r2, err)
	fmt.Printf("bufsize: %+v\n", bufSize)
	// s := string(utf16.Decode(buf))
	s := windows.UTF16PtrToString(&buf[0])
	fmt.Printf("s: %q\n", s)
}
