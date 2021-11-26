package go32

/*

BOOL GetDefaultPrinter(
  _In_    LPTSTR  pszBuffer,
  _Inout_ LPDWORD pcchBuffer
);

var(
winSpool= syscall.MustLoadDLL("Winspool.drv")
getDefaultPrinter = winSpool.MustFindProc("GetDefaultPrinterW")
)

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

*/
