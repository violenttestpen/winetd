package main

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procSuspendThread = modkernel32.NewProc("SuspendThread")

// SuspendThread suspend the execution of the main thread of the process associated with PID.
func SuspendThread(pid uint32) error {
	hSnapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(hSnapshot)

	var threadEntry windows.ThreadEntry32
	threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))

	err = windows.Thread32First(hSnapshot, &threadEntry)
	for err == nil {
		if threadEntry.OwnerProcessID == pid {
			hThread, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, threadEntry.ThreadID)
			if err != nil {
				return err
			}

			suspendThread(hThread)
			windows.CloseHandle(hThread)
			break
		}
		err = windows.Thread32Next(hSnapshot, &threadEntry)
	}
	return err
}

func suspendThread(thread windows.Handle) (uint32, error) {
	r0, _, e1 := syscall.Syscall(procSuspendThread.Addr(), 1, uintptr(thread), 0, 0)
	return uint32(r0), e1
}

// ResumeThread resumes the execution of the main thread of the process associated with PID.
func ResumeThread(pid uint32) error {
	hSnapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(hSnapshot)

	var threadEntry windows.ThreadEntry32
	threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))

	err = windows.Thread32First(hSnapshot, &threadEntry)
	for err == nil {
		if threadEntry.OwnerProcessID == pid {
			hThread, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, threadEntry.ThreadID)
			if err != nil {
				return err
			}

			windows.ResumeThread(hThread)
			windows.CloseHandle(hThread)
			break
		}
		err = windows.Thread32Next(hSnapshot, &threadEntry)
	}
	return err
}
