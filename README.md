# Windows Internet Service Daemon

`winetd` is an open-source super-server daemon that aims to replicate the functionality of `xinetd` running on Microsoft Windows.

# Installation

```
go install github.com/violenttestpen/winetd@latest
```

# Usage

```
Usage of winetd.exe:
  -bind string
        Address to bind to (default "0.0.0.0")
  -integrity string
        Run service with assigned integrity level: [Untrusted Low] (default "Untrusted")
  -port int
        Port number to bind to (default 8080)
  -server string
        Path to service to be daemonized
  -timeout int
        Timeout in seconds before closing an inactive connection (default 30)
  -verbosity int
        Verbosity mode (0-2)
```