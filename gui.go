package gui

import (
    "fmt"
    "syscall"
    "unsafe"
)

//////////////////////////////////////////////////////////////
// Guitest
// Implemented in ../runtime/syscall_windows.go.
// func Syscall(trap, nargs, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
// func Syscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
// func Syscall9(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r1, r2 uintptr, err Errno)
// func Syscall12(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12 uintptr) ...
// func Syscall15(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) ...

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
)

// POINT pos
type POINT struct {
    X, Y int32
}

// RECT window rect
type RECT struct {
    Left, Top, Right, Bottom int32
}

// HWND windows window handle
type HWND uintptr

var (
    mod = syscall.NewLazyDLL("user32.dll")

    getCursorPos     = mod.NewProc("GetCursorPos")
    setCursorPos     = mod.NewProc("SetCursorPos")

    getDesktopWindow = mod.NewProc("GetDesktopWindow")
    findWindow       = mod.NewProc("FindWindowW")
    enumChildWindows = mod.NewProc("EnumChildWindows")
    
    getWindowRect    = mod.NewProc("GetWindowRect")
    setWindowPos     = mod.NewProc("SetWindowPos")
    // getWindowLong    = mod.NewProc("GetWindowLong") // x
    
    mouse_event      = mod.NewProc("mouse_event")
    
    getClassName     = mod.NewProc("GetClassNameW")
    getWindowText    = mod.NewProc("GetWindowTextW")

    // getKeyState      = mod.NewProc("GetKeyState")
    // getKeyboardState = mod.NewProc("GetKeyboardState")
    getAsyncKeyState = mod.NewProc("GetAsyncKeyState")

    sendMessage      = mod.NewProc("SendMessageW")

    messageBox       = mod.NewProc("MessageBoxW")
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

// MessageBox popup messagebox
func MessageBox(hwnd HWND, text, caption string, btype uint32) {
    ptext    := unsafe.Pointer(syscall.StringToUTF16Ptr(text))
    pcaption := unsafe.Pointer(syscall.StringToUTF16Ptr(caption))
    syscall.Syscall6(messageBox.Addr(), 4, uintptr(hwnd), uintptr(ptext), uintptr(pcaption), uintptr(btype), 0, 0)
}

// SendMessage send message
func SendMessage(hwnd HWND, msg ,wparam, lparam uint32) {
    syscall.Syscall6(sendMessage.Addr(), 4, uintptr(hwnd), uintptr(msg),uintptr(wparam), uintptr(lparam), 0, 0)
}

// SetWindowText send wm_settext message
func SetWindowText(hwnd HWND, text string) {
    ptext    := unsafe.Pointer(syscall.StringToUTF16Ptr(text))
    syscall.Syscall6(sendMessage.Addr(), 4, uintptr(hwnd), uintptr(WM_SETTEXT),uintptr(0),uintptr(ptext), 0, 0)
}

// CloseWindow close window
func CloseWindow(hwnd HWND) {
    SendMessage(hwnd, WM_CLOSE, 0, 0)
}

// GetCursorPos get cursor pos
func GetCursorPos() POINT {
     p := POINT{}
     getCursorPos.Call(uintptr(unsafe.Pointer(&p)))
     return p
}

// SetCursorPos set cursor pos
func SetCursorPos(x, y int32) {
     syscall.Syscall(setCursorPos.Addr(), 2, uintptr(x), uintptr(y), 0)
}

// GetDesktopWindow get desktop window
func GetDesktopWindow() HWND {
     r1, _, _ := syscall.Syscall(getDesktopWindow.Addr(), 0, 0, 0, 0)
     return HWND(r1)
}

// FindWindow find className and windowName
//  ex: wnd := FindWindow("", "testwindow")
func FindWindow(className, windowName string) (w HWND, err error) {
    var cn unsafe.Pointer
    var wn unsafe.Pointer
    if len(className) > 0 {
        cn = unsafe.Pointer(syscall.StringToUTF16Ptr(className))
    }
    if len(windowName) > 0 {
        wn = unsafe.Pointer(syscall.StringToUTF16Ptr(windowName))
    }
   
    // hwnd, _, _ := findWindow.Call(uintptr(cn), uintptr(wn), 0)
    // w = HWND(hwnd)
	r0, _, e1 := syscall.Syscall(findWindow.Addr(),
        2, uintptr(cn), uintptr(wn), 0)
	w = HWND(r0)
	if w == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}

    return
}

// GetWindowText get window text(title)
func GetWindowText(h HWND) string {
    const bufSiz = 512
    var buf [bufSiz]uint16

    siz, _, _ := getWindowText.Call(uintptr(h), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
    if siz == 0 {
        return ""
    }
    name := syscall.UTF16ToString(buf[:siz])
    if siz == bufSiz-1 {
        name = name + "\u22EF"
    }
    return name
}

// GetClassName get class name
func GetClassName(h HWND) string {
    const bufSiz = 512
    var buf [bufSiz]uint16

    siz, _, _ := getClassName.Call(uintptr(h), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
    if siz == 0 {
        return ""
    }
    name := syscall.UTF16ToString(buf[:siz])
    if siz == bufSiz-1 {
        name = name + "\u22EF"
    }
    return name
}

// GetChildWindows get child windows
func GetChildWindows(h HWND) []HWND {
    hwnds := []HWND{};
    f := func(h HWND, lparam uintptr) int {
        // name := GetWindowText(h)
        // fmt.Println("NNNN: ", name)
        hwnds = append(hwnds, h)
        return 1
    }
    enumChildWindows.Call(uintptr(h), syscall.NewCallback(f), 0)
    return hwnds
}

// GetWindowRect get window rect
func GetWindowRect(h HWND) RECT {
    r := RECT{}
    getWindowRect.Call(uintptr(h),
        uintptr(unsafe.Pointer(&r.Left)),
        uintptr(unsafe.Pointer(&r.Top)),
        uintptr(unsafe.Pointer(&r.Right)),
        uintptr(unsafe.Pointer(&r.Bottom)))
    return r
}

// SetWindowPos set window poss
func SetWindowPos(h HWND, posx, posy, posex, posey int32) {
    width  := posex - posx
    height := posey - posy
    syscall.Syscall9(setWindowPos.Addr(), 7,
        uintptr(h),
        0,
        uintptr(posx),
        uintptr(posy),
        uintptr(width),
        uintptr(height),
        0,
        0,0)
}

// SetWindowSize set window pos and size
func SetWindowSize(h HWND, posx, posy, width, height int32) {
    syscall.Syscall9(setWindowPos.Addr(), 7,
        uintptr(h),
        0,
        uintptr(posx),
        uintptr(posy),
        uintptr(width),
        uintptr(height),
        0,
        0,0)
}

// --> panic: Failed to find GetWindowLong procedure in user32.dll: The specified procedure could not be found.
// const (
//     GWL_WNDPROC   = -4  // ウィンドウプロシージャのアドレスまたはウィンドウプロシージャのアドレスを示すハンドルを取得します
//     GWL_HINSTANCE = -6  // アプリケーションのインスタンスハンドルを取得します
//     GWL_HWNDPARENT= -8  // アプリケーションのインスタンスハンドルを取得します
//     GWL_ID        = -12 // ウィンドウの ID を取得します
//     GWL_STYLE     = -16 // ウィンドウスタイルを取得します
//     GWL_EXSTYLE   = -20 // 拡張ウィンドウスタイルを取得します
// )
// // window information (nindex: GWL_xxxx...)
// func GetWindowLong(h HWND, nindex int) int32 {
//     r1, _, _ := syscall.Syscall(getWindowLong.Addr(), 2,
//         uintptr(h),
//         uintptr(nindex),
//         0)
//     return int32(r1)
// }
    
const (
    // LEFTDOWN mouse button
    LEFTDOWN   = 0x00000002
    // LEFTUP mouse button
    LEFTUP     = 0x00000004
    // MIDDLEDOWN mouse button
    MIDDLEDOWN = 0x00000020
    // MIDDLEUP mouse button
    MIDDLEUP   = 0x00000040
    // MOVE mouse button
    MOVE       = 0x00000001
    // ABSOLUTE mouse button
    ABSOLUTE   = 0x00008000
    // RIGHTDOWN mouse button
    RIGHTDOWN  = 0x00000008
    // RIGHTUP mouse button
    RIGHTUP    = 0x00000010
    // WHEEL mouse button
    WHEEL      = 0x00000800
    // HWHEEL mouse button
    HWHEEL     = 0x00001000
)

// MouseMoveAbs mouse move absolute
func MouseMoveAbs(x, y int) {
    bit := MOVE | ABSOLUTE
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        uintptr(x),uintptr(y),
        0,0,0)
}

// MouseMoveRel mouse move relocate
func MouseMoveRel(x, y int) {
    bit := MOVE
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        uintptr(x),uintptr(y),
        0,0,0)
}

// MouseLButtonDown mouse button down(left)
func MouseLButtonDown() {
    bit := LEFTDOWN
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// MouseLButtonUp mouse button up(left)
func MouseLButtonUp() {
    bit := LEFTUP
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// MouseRButtonDown mouse button down(right)
func MouseRButtonDown() {
    bit := RIGHTDOWN
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// MouseRButtonUp mouse button up(right)
func MouseRButtonUp() {
    bit := RIGHTUP
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// MouseMButtonDown mouse button down(middle)
func MouseMButtonDown() {
    bit := MIDDLEDOWN
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// MouseMButtonUp mouse button up(middle)
func MouseMButtonUp() {
    bit := MIDDLEUP
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// ClickMouseLeft click pos left button
func ClickMouseLeft(posx, posy int32) {
    SetCursorPos(posx, posy)
    bit := LEFTDOWN | LEFTUP
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// ClickMouseRight click pos right button
func ClickMouseRight(posx, posy int32) {
    SetCursorPos(posx, posy)
    bit := RIGHTDOWN | RIGHTUP
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// ClickMouseMiddle click pos middle button
func ClickMouseMiddle(posx, posy int32) {
    SetCursorPos(posx, posy)
    bit := MIDDLEDOWN | MIDDLEUP
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,0,0,
        0)
}

// MouseMoveWheel move mouse wheel(side)
func MouseMoveWheel(scroll int32) {
    bit := WHEEL
    scroll *= WHEEL_DELTA
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,
        uintptr(scroll),
        0,
        0)
}

// MouseMoveHWheel move mouse wheel(height)
func MouseMoveHWheel(scroll int32) {
    bit := HWHEEL
    scroll *= WHEEL_DELTA
    syscall.Syscall6(mouse_event.Addr(), 5,
        uintptr(bit),
        0,0,
        uintptr(scroll),
        0,
        0)
}

// func GetKeyState(vkey int) int {
//     r1, _, _ := getKeyState.Call(uintptr(vkey))
//     return int(r1)
// }
// func GetKeyboardState(vkey uint8) {
//  const bufSiz = 256
//     var buf [bufSiz]byte
//     r1, _, e1 := syscall.Syscall(getKeyboardState.Addr(), 1, uintptr(unsafe.Pointer(&buf)), 0, 0)
// }

// GetAsyncKeyState virtual key
func GetAsyncKeyState(vkey byte) int {
    r1, _, _ := getAsyncKeyState.Call(uintptr(vkey))
    return int(uint16(r1))
}

var keymaps map[string]byte = map[string]byte{
    "LBUTTON" : 0x01,
    "RBUTTON" : 0x02,
    "CANCEL"  : 0x03,
    "MBUTTON" : 0x04,
    "XBUTTON1": 0x05,
    "XBUTTON2": 0x06,
    "BACK"    : 0x08,
    "TAB"     : 0x09,
    "RETURN"  : 0x0d,
    "SHIFT"   : 0x10,
    "CTRL"    : 0x11,
    "PAUSE"   : 0x13,
    "CAPSLOCK": 0x14,
    
    "ESC"     : 0x1b,
    "SPACE"   : 0x20,
    "PAGEUP"  : 0x21,
    "PAGEDOWN": 0x22,
    "END"     : 0x23,
    "HOME"    : 0x24,
    "LEFT"    : 0x25,
    "UP"      : 0x26,
    "RIGHT"   : 0x27,
    "DOWN"    : 0x28,

    "PRINT"   : 0x2a,
    "INSERT"  : 0x2d,
    "DELETE"  : 0x2e,
    "HELP"    : 0x2f,
    
    " "       : ' ',
    "A" : 'A',
    "B" : 'B',
    "C" : 'C',
    "D" : 'D',
    "E" : 'E',
    "F" : 'F',
    "G" : 'G',
    "H" : 'H',
    "I" : 'I',
    "J" : 'J',
    "K" : 'K',
    "L" : 'L',
    "M" : 'M',
    "N" : 'N',
    "O" : 'O',
    "P" : 'P',
    "Q" : 'Q',
    "R" : 'R',
    "S" : 'S',
    "T" : 'T',
    "U" : 'U',
    "V" : 'V',
    "W" : 'W',
    "X" : 'X',
    "Y" : 'W',
    "Z" : 'Z',
    "a" : 'A',
    "b" : 'B',
    "c" : 'C',
    "d" : 'D',
    "e" : 'E',
    "f" : 'F',
    "g" : 'G',
    "h" : 'H',
    "i" : 'I',
    "j" : 'J',
    "k" : 'K',
    "l" : 'L',
    "m" : 'M',
    "n" : 'N',
    "o" : 'O',
    "p" : 'P',
    "q" : 'Q',
    "r" : 'R',
    "s" : 'S',
    "t" : 'T',
    "u" : 'U',
    "v" : 'V',
    "w" : 'W',
    "x" : 'X',
    "y" : 'W',
    "z" : 'Z',

    "LWNN": 0x5b,
    "RWNN": 0x5c,
    "APP" : 0x5d, // menu
   
    // tenkey
    "n0" : 0x60,
    "n1" : 0x61,
    "n2" : 0x62,
    "n3" : 0x63,
    "n4" : 0x64,
    "n5" : 0x65,
    "n6" : 0x66,
    "n7" : 0x67,
    "n8" : 0x68,
    "n9" : 0x69,
    "n*"   : 0x6a, // tenkey * 
    "MULTIPLY": 0x6a, // tenkey * 
    "n+"   : 0x6b,   // tenkey +
    "ADD" : 0x6b,
    // "ENTER": 0x6c,  // tenkey enter 来ない
    "SUBTRACT": 0x6d,
    "n-"   :  0x6d,
    "DECIMAL": 0x6e, // tenkey .
    "DOT" :    0x6e, // tenkey .
    "n." :    0x6e, // tenkey .
    "DIVIDE":  0x6f, // tenkey /
    "n/"     :  0x6f, // tenkey /

    "F1"   : 0x70,
    "F2"   : 0x71,
    "F3"   : 0x72,
    "F4"   : 0x73,
    "F5"   : 0x74,
    "F6"   : 0x75,
    "F7"   : 0x76,
    "F8"   : 0x77,
    "F9"   : 0x78,
    "F10"   : 0x79,
    "F11"   : 0x7a,
    "F12"   : 0x7b,
    "F13"   : 0x7c,
    "F14"   : 0x7d,
    "F15"   : 0x7e,
    "F16"   : 0x7f,
    "F17"   : 0x80,
    "F18"   : 0x81,
    "F19"   : 0x82,
    "F20"   : 0x83,
    "F21"   : 0x84,
    "F22"   : 0x85,
    "F23"   : 0x86,
    "F24"   : 0x87,

    "NUMLOCK": 0x90,
    "SCROLL":  0x91,

    "LSHIFT" : 0xa0,
    "RSHIFT" : 0xa1,
    "LCTRL"  : 0xa2,
    "RCTRL"  : 0xa3,
    "LMENU"  : 0xa4,
    "RMENU"  : 0xa5,

    "BROWSER_BACK" : 0xA6,
    "BROWSER_FORWARD": 0xA7,
    "BROWSER_REFRESH": 0xA8,
    "BROWSER_STOP": 0xA9,
    "BROWSER_SEARCH": 0xAA,
    "BROWSER_FAVORITES": 0xAB,
    "BROWSER_HOME": 0xAC,
    "VOLUME_MUTE":  0xAD,
    "VOLUME_DOWN":  0xAE,
    "VOLUME_UP":    0xAF,

    "MEDIA_NEXT_TRACK": 0xB0,
    "MEDIA_PREV_TRACK": 0xB1,
    "MEDIA_STOP":   0xB2,
    "MEDIA_PLAY_PAUSE": 0xB3,
    "LAUNCH_MAIL":  0xB4,
    "LAUNCH_MEDIA_SELECT":  0xB5,
    "LAUNCH_APP1":  0xB6,
    "LAUNCH_APP2":  0xB7,
    
    ":" : 0xba,
    "*" : 0xba,
    ";" : 0xbb,
    "+" : 0xbb,
    "," : 0xbc,
    "." : 0xbe,
    "/" : 0xbf,
    "?" : 0xbf,
    "`" : 0xc0,
    "@" : 0xc0,
    
    "{" : 0xdb,
    "[" : 0xdb,
    "\\": 0xdc,
    "|" : 0xdc,
    "}" : 0xdd,
    "]" : 0xdd,
    "'" : 0xde,
    "\"" : 0xde,
    
    "ZOOM": 0xfb,    
}

// PressedKey key check pressed
func PressedKey(key string) bool {
    v, ok := keymaps[key]
    if ok {
        return (GetAsyncKeyState(v) != 0)
    }
    fmt.Println("not found key name: ", key)
    return  false
}

// IsValidKey valid key
func IsValidKey(key string) (ok bool) {
    _, ok = keymaps[key]
    return
}

// KeyList valid key list
func KeyList() (list []string) {
    for k, _ := range keymaps {
        list = append(list, k)
    }
    return
}

func printTreeClassNameSub(hwnd HWND, nest int) {
    space := ""
    for n := 0; n < nest; n++ {
        space = space + "  "
    }
    str := GetClassName(hwnd)
    fmt.Println(space + str)

    hwnds := GetChildWindows(hwnd)
    for _, chwnd := range hwnds {
        printTreeClassNameSub(chwnd, nest + 1)
    }
}

// PrintTreeClassName print class name
func PrintTreeClassName(hwnd HWND) {
    // hwnd = GetDesktopWindow()
    fmt.Println("--- PrintTreeClassName ---")
    printTreeClassNameSub(hwnd, 0)
    fmt.Println("--------------------------")
}
