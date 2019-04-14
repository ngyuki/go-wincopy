package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

func MustAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return abs
}

func copyToClipboard(files []string) {

	// UTF16 で {$file}\0{$file}\0{$file}\0 な形式のバッファ
	var ufiles []uint16
	for _, v := range files {
		a := MustAbsPath(v)
		fmt.Printf("copy: %v\n", a)
		f := syscall.StringToUTF16(MustAbsPath(v))
		ufiles = append(ufiles, f...)
	}

	// グローバルメモリサイズ ... DROPFILES + $files + \0
	size := uint32(unsafe.Sizeof(DROPFILES{})) + uint32(len(ufiles)*2) + 2

	mem := GlobalAlloc(size)
	defer GlobalFree(mem)

	{
		buf := (*DROPFILES)(GlobalLock(mem))
		defer GlobalUnlock(mem)

		buf.pFiles = uint32(unsafe.Sizeof(DROPFILES{}))
		buf.ptX = 0
		buf.ptY = 0
		buf.fNC = 0
		buf.fWide = 1

		// 構造体の終端のポインタに files をコピー
		pfiles := (*[4096]uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(buf)) + unsafe.Sizeof(DROPFILES{})))
		copy(pfiles[:], ufiles)
	}

	OpenClipboard()
	defer CloseClipboard()
	EmptyClipboard()
	SetClipboardData(CF_HDROP, mem)
}

func getFromClipboard() []string {

	OpenClipboard()
	defer CloseClipboard()

	mem := GetClipboardData(CF_HDROP)

	buf := (*DROPFILES)(GlobalLock(mem))
	defer GlobalUnlock(mem)

	var files []string

	ptr := (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(buf)) + uintptr(buf.pFiles)))
	for *ptr != 0 {
		var buf []uint16
		size := unsafe.Sizeof(uint16(0))
		for *ptr != 0 {
			buf = append(buf, *ptr)
			ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + size))
		}
		file := syscall.UTF16ToString(buf)
		if len(file) > 0 {
			files = append(files, file)
		}
		ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + size))
	}

	return files
}

func copyFile(dstName string, srcName string) {

	src, err := os.Open(srcName)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
}

func pasteToSave(files []string) {
	for _, file := range files {
		src := MustAbsPath(file)
		dst := MustAbsPath(filepath.Base(src))
		if dst == src {
			fmt.Printf("skip: %v\n", src)
		} else {
			fmt.Printf("copy: %v\n   -> %v\n", src, dst)
			copyFile(dst, src)
		}
	}
}

func main() {
	paste := flag.Bool("p", false, "paste flag")
	flag.Parse()

	if *paste {
		files := getFromClipboard()
		pasteToSave(files)
	} else {
		copyToClipboard(flag.Args())
	}
}
