// Code generated by 'go generate'; DO NOT EDIT.

package windows

import (
	"syscall"
	"unsafe"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")
	modiphlpapi = syscall.NewLazyDLL("iphlpapi.dll")
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	modnetapi32 = syscall.NewLazyDLL("netapi32.dll")
	modpsapi    = syscall.NewLazyDLL("psapi.dll")
	moduserenv  = syscall.NewLazyDLL("userenv.dll")
	modws2_32   = syscall.NewLazyDLL("ws2_32.dll")

	procAdjustTokenPrivileges = modadvapi32.NewProc("AdjustTokenPrivileges")
	procDuplicateTokenEx      = modadvapi32.NewProc("DuplicateTokenEx")
	procImpersonateSelf       = modadvapi32.NewProc("ImpersonateSelf")
	procLookupPrivilegeValueW = modadvapi32.NewProc("LookupPrivilegeValueW")
	procOpenThreadToken       = modadvapi32.NewProc("OpenThreadToken")
	procRevertToSelf          = modadvapi32.NewProc("RevertToSelf")
	procSetTokenInformation   = modadvapi32.NewProc("SetTokenInformation")
)

func adjustTokenPrivileges(token syscall.Token, disableAllPrivileges bool, newstate *TOKEN_PRIVILEGES, buflen uint32, prevstate *TOKEN_PRIVILEGES, returnlen *uint32) (ret uint32, err error) {
	var _p0 uint32
	if disableAllPrivileges {
		_p0 = 1
	}
	r0, _, e1 := syscall.Syscall6(procAdjustTokenPrivileges.Addr(), 6, uintptr(token), uintptr(_p0), uintptr(unsafe.Pointer(newstate)), uintptr(buflen), uintptr(unsafe.Pointer(prevstate)), uintptr(unsafe.Pointer(returnlen)))
	ret = uint32(r0)
	if true {
		err = errnoErr(e1)
	}
	return
}

func DuplicateTokenEx(hExistingToken syscall.Token, dwDesiredAccess uint32, lpTokenAttributes *syscall.SecurityAttributes, impersonationLevel uint32, tokenType TokenType, phNewToken *syscall.Token) (err error) {
	r1, _, e1 := syscall.Syscall6(procDuplicateTokenEx.Addr(), 6, uintptr(hExistingToken), uintptr(dwDesiredAccess), uintptr(unsafe.Pointer(lpTokenAttributes)), uintptr(impersonationLevel), uintptr(tokenType), uintptr(unsafe.Pointer(phNewToken)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func ImpersonateSelf(impersonationlevel uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procImpersonateSelf.Addr(), 1, uintptr(impersonationlevel), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func LookupPrivilegeValue(systemname *uint16, name *uint16, luid *LUID) (err error) {
	r1, _, e1 := syscall.Syscall(procLookupPrivilegeValueW.Addr(), 3, uintptr(unsafe.Pointer(systemname)), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(luid)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func OpenThreadToken(h syscall.Handle, access uint32, openasself bool, token *syscall.Token) (err error) {
	var _p0 uint32
	if openasself {
		_p0 = 1
	}
	r1, _, e1 := syscall.Syscall6(procOpenThreadToken.Addr(), 4, uintptr(h), uintptr(access), uintptr(_p0), uintptr(unsafe.Pointer(token)), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func RevertToSelf() (err error) {
	r1, _, e1 := syscall.Syscall(procRevertToSelf.Addr(), 0, 0, 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetTokenInformation(tokenHandle syscall.Token, tokenInformationClass uint32, tokenInformation uintptr, tokenInformationLength uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procSetTokenInformation.Addr(), 4, uintptr(tokenHandle), uintptr(tokenInformationClass), uintptr(tokenInformation), uintptr(tokenInformationLength), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}