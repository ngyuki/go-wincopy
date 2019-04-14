package main

import (
	"syscall"
	"unsafe"
)

type DROPFILES struct {
	pFiles uint32 // DWORD
	ptX    int32  // LONG
	ptY    int32  // LONG
	fNC    int32  // BOOL
	fWide  int32  // BOOL
}

const CF_HDROP = 15

func GlobalAlloc(size uint32) uintptr {
	dll := syscall.MustLoadDLL("Kernel32.dll")
	defer dll.Release()
	const GMEM_MOVEABLE = 0x0002
	const GMEM_ZEROINIT = 0x0040
	r, _, _ := dll.MustFindProc("GlobalAlloc").Call(GMEM_MOVEABLE|GMEM_ZEROINIT, uintptr(size))
	if r == 0 {
		err := syscall.GetLastError()
		panic(err)
	}
	return r
}

func GlobalLock(mem uintptr) unsafe.Pointer {
	dll := syscall.MustLoadDLL("Kernel32.dll")
	defer dll.Release()
	r, _, _ := dll.MustFindProc("GlobalLock").Call(mem)
	if r == 0 {
		err := syscall.GetLastError()
		panic(err)
	}
	return unsafe.Pointer(r)
}

func GlobalUnlock(mem uintptr) {
	dll := syscall.MustLoadDLL("Kernel32.dll")
	defer dll.Release()
	r, _, _ := dll.MustFindProc("GlobalUnlock").Call(mem)
	if r == 0 {
		err := syscall.GetLastError()
		if err != nil {
			panic(err)
		}
	}
}

func GlobalFree(mem uintptr) {
	dll := syscall.MustLoadDLL("Kernel32.dll")
	defer dll.Release()
	r, _, _ := dll.MustFindProc("GlobalFree").Call(mem)
	if r == 0 {
		err := syscall.GetLastError()
		if err != nil {
			panic(err)
		}
	}
}

func OpenClipboard() {
	dll := syscall.MustLoadDLL("user32.dll")
	defer dll.Release()
	r, _, _ := dll.MustFindProc("OpenClipboard").Call(0)
	if r == 0 {
		err := syscall.GetLastError()
		panic(err)
	}
}

func GetClipboardData(format uint32) uintptr {
	dll := syscall.MustLoadDLL("user32.dll")
	defer dll.Release()
	r, _, _ := dll.MustFindProc("GetClipboardData").Call(uintptr(format))
	if r == 0 {
		err := syscall.GetLastError()
		panic(err)
	}
	return r
}

func SetClipboardData(format uint32, mem uintptr) uintptr {
	dll := syscall.MustLoadDLL("user32.dll")
	defer dll.Release()
	r, _, _ := dll.MustFindProc("SetClipboardData").Call(uintptr(format), mem)
	if r == 0 {
		err := syscall.GetLastError()
		panic(err)
	}
	return r
}

func EmptyClipboard() {
	dll := syscall.MustLoadDLL("user32.dll")
	defer dll.Release()
	r, _, _ := dll.MustFindProc("EmptyClipboard").Call()
	if r == 0 {
		err := syscall.GetLastError()
		panic(err)
	}
}

func CloseClipboard() {
	dll := syscall.MustLoadDLL("user32.dll")
	defer dll.Release()
	r, _, _ := dll.MustFindProc("CloseClipboard").Call(0)
	if r == 0 {
		err := syscall.GetLastError()
		panic(err)
	}
}
