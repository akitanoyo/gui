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

type POINT struct {
	X, Y int32
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

type HWND uintptr

var (
    mod = syscall.NewLazyDLL("user32.dll")

    getCursorPos     = mod.NewProc("GetCursorPos")
    setCursorPos     = mod.NewProc("SetCursorPos")

    getDesktopWindow = mod.NewProc("GetDesktopWindow")
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

func MessageBox(hwnd HWND, text, caption string, btype uint32) {
    ptext    := unsafe.Pointer(syscall.StringToUTF16Ptr(text))
    pcaption := unsafe.Pointer(syscall.StringToUTF16Ptr(caption))
    syscall.Syscall6(messageBox.Addr(), 4, uintptr(hwnd), uintptr(ptext), uintptr(pcaption), uintptr(btype), 0, 0)
}

func SendMessage(hwnd HWND, msg ,wparam, lparam uint32) {
    syscall.Syscall6(sendMessage.Addr(), 4, uintptr(hwnd), uintptr(msg),uintptr(wparam), uintptr(lparam), 0, 0)
}

func CloseWindow(hwnd HWND) {
	SendMessage(hwnd, WM_CLOSE, 0, 0)
}

func GetCursorPos() POINT {
     p := POINT{}
     getCursorPos.Call(uintptr(unsafe.Pointer(&p)))
     return p
}

func SetCursorPos(x, y int32) {
     syscall.Syscall(setCursorPos.Addr(), 2, uintptr(x), uintptr(y), 0)
}

func GetDesktopWindow() HWND {
     r1, _, _ := syscall.Syscall(getDesktopWindow.Addr(), 0, 0, 0, 0)
     return HWND(r1)
}

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

func GetWindowRect(h HWND) RECT {
	r := RECT{}
	getWindowRect.Call(uintptr(h),
		uintptr(unsafe.Pointer(&r.Left)),
		uintptr(unsafe.Pointer(&r.Top)),
		uintptr(unsafe.Pointer(&r.Right)),
		uintptr(unsafe.Pointer(&r.Bottom)))
	return r
}

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
	LEFTDOWN   = 0x00000002
	LEFTUP     = 0x00000004
	MIDDLEDOWN = 0x00000020
	MIDDLEUP   = 0x00000040
	MOVE       = 0x00000001
	ABSOLUTE   = 0x00008000
	RIGHTDOWN  = 0x00000008
	RIGHTUP    = 0x00000010
	WHEEL      = 0x00000800
    HWHEEL     = 0x00001000
)

func MouseMoveAbs(x, y int) {
	bit := MOVE | ABSOLUTE
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		uintptr(x),uintptr(y),
        0,0,0)
}

func MouseMoveRel(x, y int) {
	bit := MOVE
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		uintptr(x),uintptr(y),
		0,0,0)
}

func MouseLButtonDown() {
	bit := LEFTDOWN
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

func MouseLButtonUp() {
	bit := LEFTUP
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

func MouseRButtonDown() {
	bit := RIGHTDOWN
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

func MouseRButtonUp() {
	bit := RIGHTUP
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

func MouseMButtonDown() {
	bit := MIDDLEDOWN
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

func MouseMButtonUp() {
	bit := MIDDLEUP
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

func ClickMouseLeft(posx, posy int32) {
	SetCursorPos(posx, posy)
	bit := LEFTDOWN | LEFTUP
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

func ClickMouseRight(posx, posy int32) {
	SetCursorPos(posx, posy)
	bit := RIGHTDOWN | RIGHTUP
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

func ClickMouseMiddle(posx, posy int32) {
	SetCursorPos(posx, posy)
	bit := MIDDLEDOWN | MIDDLEUP
	syscall.Syscall6(mouse_event.Addr(), 5,
		uintptr(bit),
		0,0,0,0,
		0)
}

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
// 	const bufSiz = 256
//     var buf [bufSiz]byte
//     r1, _, e1 := syscall.Syscall(getKeyboardState.Addr(), 1, uintptr(unsafe.Pointer(&buf)), 0, 0)
// }

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
    "BROWSER_STOP":	0xA9,
    "BROWSER_SEARCH": 0xAA,
    "BROWSER_FAVORITES": 0xAB,
    "BROWSER_HOME":	0xAC,
    "VOLUME_MUTE":	0xAD,
    "VOLUME_DOWN":	0xAE,
    "VOLUME_UP":	0xAF,

    "MEDIA_NEXT_TRACK":	0xB0,
    "MEDIA_PREV_TRACK":	0xB1,
    "MEDIA_STOP":	0xB2,
    "MEDIA_PLAY_PAUSE":	0xB3,
    "LAUNCH_MAIL":	0xB4,
    "LAUNCH_MEDIA_SELECT":	0xB5,
    "LAUNCH_APP1":	0xB6,
    "LAUNCH_APP2":	0xB7,
    
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

func PressedKey(key string) bool {
    v, ok := keymaps[key]
    if ok {
        return (GetAsyncKeyState(v) != 0)
    }
    fmt.Println("not found key name: ", key)
    return  false
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

func PrintTreeClassName(hwnd HWND) {
    // hwnd = GetDesktopWindow()
    fmt.Println("--- PrintTreeClassName ---")
    printTreeClassNameSub(hwnd, 0)
    fmt.Println("--------------------------")
}
