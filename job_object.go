package main

import (
	"os"
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"
)

// NewJobFromProcess assigns the process to a new job object
func NewJobFromProcess(p *os.Process) (windows.Handle, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return windows.InvalidHandle, err
	}

	processValue := reflect.ValueOf(p).Elem()
	hProcess := processValue.FieldByName("handle").Uint()
	if err := windows.AssignProcessToJobObject(job, windows.Handle(hProcess)); err != nil {
		return windows.InvalidHandle, err
	}

	return job, nil
}

// ResumeThread resumes the execution of the main thread of the process associated with PID
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

// TerminateJob terminates the job object and its child processes before closing the handle
func TerminateJob(job windows.Handle) error {
	if err := windows.TerminateJobObject(windows.Handle(job), 0); err != nil {
		return err
	}
	return windows.CloseHandle(job)
}
