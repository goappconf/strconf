package strconf

import (
	"os/exec"
	"testing"
)

func TestHideWindow(t *testing.T) {
	cmd := exec.Command("cmd.exe")
	hideWindow(cmd)
	if cmd.SysProcAttr == nil || !cmd.SysProcAttr.HideWindow || cmd.SysProcAttr.CreationFlags&0x08000000 == 0 {
		t.Fatal("expected hidden window and CREATE_NO_WINDOW")
	}
}
