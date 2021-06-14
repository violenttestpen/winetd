package main

import (
	"syscall"
	"unsafe"

	"github.com/violenttestpen/ginetd/internal/syscall/windows"
)

var sidWinIntegrityLevels = map[string]string{
	"Untrusted": `S-1-16-0`,
	"Low":       `S-1-16-4096`,
}

func getIntegrityLevelToken(wns string) (syscall.Token, error) {
	var procToken, token syscall.Token
	proc, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, err
	}
	defer syscall.CloseHandle(proc)

	var access uint32 = syscall.TOKEN_DUPLICATE | syscall.TOKEN_ADJUST_DEFAULT |
		syscall.TOKEN_QUERY | syscall.TOKEN_ASSIGN_PRIMARY
	err = syscall.OpenProcessToken(proc, access, &procToken)
	if err != nil {
		return 0, err
	}
	defer procToken.Close()

	err = windows.DuplicateTokenEx(procToken, 0, nil, windows.SecurityImpersonation,
		windows.TokenPrimary, &token)
	if err != nil {
		return 0, err
	}

	sid, err := syscall.StringToSid(wns)
	if err != nil {
		return 0, err
	}

	tml := &windows.TOKEN_MANDATORY_LABEL{}
	tml.Label.Attributes = windows.SE_GROUP_INTEGRITY
	tml.Label.Sid = sid

	err = windows.SetTokenInformation(token, syscall.TokenIntegrityLevel,
		uintptr(unsafe.Pointer(tml)), tml.Size())
	if err != nil {
		token.Close()
		return 0, err
	}

	return token, nil
}
