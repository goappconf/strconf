//go:build !windows

package strconf

import "os/exec"

func hideWindow(cmd *exec.Cmd) {}
