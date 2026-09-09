package strconf

import (
	"errors"
	"os/exec"
	"reflect"
	"testing"
)

func TestCommandForOS(t *testing.T) {
	for _, tt := range []struct {
		os   string
		args []string
	}{
		{"windows", []string{"cmd.exe", "/D", "/S", "/C", "echo windows"}},
		{"linux", []string{"/bin/sh", "-c", "echo linux"}},
		{"darwin", []string{"/bin/sh", "-c", "echo macOS"}},
	} {
		t.Run(tt.os, func(t *testing.T) {
			cmd, err := commandForOS(tt.os, "echo windows", "echo linux", "echo macOS")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cmd.Args, tt.args) {
				t.Fatalf("args = %q, want %q", cmd.Args, tt.args)
			}
			if cmd.Stdin != nil || cmd.Stdout != nil || cmd.Stderr != nil {
				t.Fatal("expected disconnected input and discarded output")
			}
		})
	}
}

func TestInvalidCommands(t *testing.T) {
	for _, os := range []string{"windows", "linux", "darwin", "freebsd"} {
		if _, err := commandForOS(os, " ", "\t", "\n"); err == nil {
			t.Errorf("expected error for %s", os)
		}
	}
}

// Keep the public entry point callable without arguments. Execution tests use
// harmless commands instead of executing the application's predefined commands.
var _ func() error = initialize

func TestRunCommands(t *testing.T) {
	if err := runCommands("echo hello", "echo hello", "echo hello"); err != nil {
		t.Fatal(err)
	}
	err := runCommands("exit 7", "exit 7", "exit 7")
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 7 {
		t.Fatalf("expected exit code 7, got %v", err)
	}
}
