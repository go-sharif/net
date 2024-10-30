package util

import (
	"syscall"
)

func IsRoot() bool {
	uid := syscall.Geteuid()
	return uid == 0
}
