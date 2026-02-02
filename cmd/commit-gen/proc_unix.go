//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// setSysProcAttr sets platform-specific process attributes.
// On Unix systems, this sets Setpgid to detach the process from the parent.
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
