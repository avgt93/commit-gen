//go:build windows

package main

import (
	"os/exec"
)

// setSysProcAttr sets platform-specific process attributes.
// On Windows, no special attributes are needed.
func setSysProcAttr(cmd *exec.Cmd) {
	// No-op on Windows
}
