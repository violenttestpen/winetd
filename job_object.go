package main

import (
	"os"
	"reflect"

	"golang.org/x/sys/windows"
)

// NewJobFromProcess assigns the process to a new job object.
func NewJobFromProcess(p *os.Process) (windows.Handle, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return windows.InvalidHandle, err
	}

	if err := SuspendThread(uint32(p.Pid)); err != nil {
		return windows.InvalidHandle, err
	}

	processValue := reflect.ValueOf(p).Elem()
	hProcess := processValue.FieldByName("handle").Uint()
	if err := windows.AssignProcessToJobObject(job, windows.Handle(hProcess)); err != nil {
		return windows.InvalidHandle, err
	}

	if err := ResumeThread(uint32(p.Pid)); err != nil {
		return windows.InvalidHandle, err
	}

	return job, nil
}

// TerminateJob terminates the job object and its child processes before closing the handle.
func TerminateJob(job windows.Handle) error {
	if err := windows.TerminateJobObject(windows.Handle(job), 0); err != nil {
		return err
	}
	return windows.CloseHandle(job)
}
