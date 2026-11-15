package main

import "golang.org/x/sys/windows"

var (
	bcrypt = windows.NewLazySystemDLL("bcrypt.dll")

	BCryptGenRandom = bcrypt.NewProc("BCryptGenRandom")
)

const bcryptUseSystemPreferredRNG = 0x00000002
