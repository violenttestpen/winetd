package windows

import (
	"os"
	"reflect"
	"unsafe"

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

// GetBasicJobLimitInfo retrieves limit and job state information from the job object.
func GetBasicJobLimitInfo(job windows.Handle) (*windows.JOBOBJECT_BASIC_LIMIT_INFORMATION, error) {
	info := new(windows.JOBOBJECT_BASIC_LIMIT_INFORMATION)
	err := windows.QueryInformationJobObject(job,
		windows.JobObjectBasicLimitInformation,
		uintptr(unsafe.Pointer(info)),
		uint32(unsafe.Sizeof(*info)), nil)
	if err != nil {
		return nil, err
	}
	return info, err
}

// GetExtendedJobLimitInfo retrieves extended limit and job state information from the job object.
func GetExtendedJobLimitInfo(job windows.Handle) (*windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION, error) {
	info := new(windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION)
	err := windows.QueryInformationJobObject(job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(info)),
		uint32(unsafe.Sizeof(*info)), nil)
	if err != nil {
		return nil, err
	}
	return info, err
}

// TerminateJob terminates the job object and its child processes before closing the handle.
func TerminateJob(job windows.Handle) error {
	if err := windows.TerminateJobObject(windows.Handle(job), 0); err != nil {
		return err
	}
	return windows.CloseHandle(job)
}
