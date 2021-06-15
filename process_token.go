package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var sidWinIntegrityLevels = map[string]uint32{
	"Untrusted": windows.WinUntrustedLabelSid, // `S-1-16-0`
	"Low":       windows.WinLowLabelSid,       // `S-1-16-4096`
}

func getIntegrityLevelToken(sidType uint32) (windows.Token, error) {
	var procToken, token windows.Token
	proc, err := windows.GetCurrentProcess()
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(proc)

	var access uint32 = windows.TOKEN_DUPLICATE | windows.TOKEN_ADJUST_DEFAULT |
		windows.TOKEN_QUERY | windows.TOKEN_ASSIGN_PRIMARY
	err = windows.OpenProcessToken(proc, access, &procToken)
	if err != nil {
		return 0, err
	}
	defer procToken.Close()

	err = windows.DuplicateTokenEx(procToken, 0, nil, windows.SecurityImpersonation,
		windows.TokenPrimary, &token)
	if err != nil {
		return 0, err
	}

	sid, err := windows.CreateWellKnownSid(windows.WELL_KNOWN_SID_TYPE(sidType))
	if err != nil {
		return 0, err
	}

	tml := &windows.Tokenmandatorylabel{}
	tml.Label.Attributes = windows.SE_GROUP_INTEGRITY
	tml.Label.Sid = sid

	err = windows.SetTokenInformation(token, windows.TokenIntegrityLevel,
		(*byte)(unsafe.Pointer(tml)), tml.Size())
	if err != nil {
		token.Close()
		return 0, err
	}

	return token, nil
}
