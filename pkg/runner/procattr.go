//go:build !windows
// +build !windows

package runner

import "syscall"

func newProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
