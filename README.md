# Go Internet Service Daemon

`ginetd` is an open-source super-server daemon that aims to replicate the functionality of `xinetd` running on Microsoft Windows.

# Installation

```
go install github.com/violenttestpen/ginetd
```

# Usage

```
Usage of ginetd.exe:
  -bind string
        Address to bind to (default "0.0.0.0")
  -pass string
        The password of the user
  -port int
        Port number to bind to (default 8080)
  -server string
        Path to service to be daemonized
  -user string
        The user to run the service as
  -verbosity int
        Verbosity mode (0-2)
```