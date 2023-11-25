package helper

import (
	"strings"
	"syscall"
	"unsafe"
)

type ulong int32
type ulong_ptr uintptr

type PROCESSENTRY32 struct {
	dwSize              ulong
	cntUsage            ulong
	th32ProcessID       ulong
	th32DefaultHeapID   ulong_ptr
	th32ModuleID        ulong
	cntThreads          ulong
	th32ParentProcessID ulong
	pcPriClassBase      ulong
	dwFlags             ulong
	szExeFile           [260]byte
}

func getProcess() (m map[string]byte) {
	m = make(map[string]byte, 100)
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	CreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	pHandle, _, _ := CreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	if int(pHandle) == -1 {
		return
	}
	defer func() {
		CloseHandle := kernel32.NewProc("CloseHandle")
		_, _, _ = CloseHandle.Call(pHandle)
	}()
	Process32Next := kernel32.NewProc("Process32Next")
	for {
		var proc PROCESSENTRY32
		proc.dwSize = ulong(unsafe.Sizeof(proc))
		if rt, _, _ := Process32Next.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
			v := strings.TrimSpace(string(proc.szExeFile[0:]))
			var i = strings.Index(v, ".exe")
			if i > 0 {
				name := v[:i]
				name = NormalizeProcessName(name)
				m[name]++
			}
		} else {
			break
		}
	}
	return m
}

func NormalizeProcessName(src string) string {
	var dst = make([]byte, 0, len(src))
	for _, b := range []byte(src) {
		if isLetter(b) || b == '_' {
			dst = append(dst, b)
		}
	}
	return string(dst)
}

func isLetter(b byte) bool {
	if b >= 'a' && b <= 'z' {
		return true
	}
	if b >= 'A' && b <= 'Z' {
		return true
	}
	return false
}
