package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/violenttestpen/winetd/pkg/log"
	"github.com/violenttestpen/winetd/pkg/windows"
)

const integrityUsage = "Run service with assigned integrity level: %v"

var (
	bind      = "0.0.0.0"
	port      = 8080
	integrity = "Untrusted"
	timeout   = 30
	verbosity = 0

	service string
)

var logger log.Log

func main() {
	i, integrityLevels := 0, make([]string, len(windows.SidWinIntegrityLevels))
	for k := range windows.SidWinIntegrityLevels {
		integrityLevels[i] = k
		i++
	}

	flag.StringVar(&bind, "bind", bind, "Address to bind to")
	flag.IntVar(&port, "port", port, "Port number to bind to")
	flag.StringVar(&integrity, "integrity", integrity, fmt.Sprintf(integrityUsage, integrityLevels))
	flag.IntVar(&timeout, "timeout", timeout, "Timeout in seconds before closing an inactive connection")
	flag.IntVar(&verbosity, "verbosity", verbosity, "Verbosity mode (0-2)")
	flag.StringVar(&service, "server", service, "Path to service to be daemonized")
	flag.Parse()

	logger = log.NewLogger(verbosity)
	if err := run(net.JoinHostPort(bind, strconv.Itoa(port)), service, integrity, timeout); err != nil {
		logger.Fatal(err)
	}
}

func run(bind, service, integrity string, timeout int) error {
	sid, ok := windows.SidWinIntegrityLevels[integrity]
	if !ok {
		return windows.ErrInvalidIntegrityLevel
	}

	token, err := windows.GetIntegrityLevelToken(sid)
	if err != nil {
		return err
	}
	defer token.Close()

	if _, err := os.Stat(service); err != nil {
		return err
	}

	server, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}
	defer server.Close()

	return doListen(server, func(conn net.Conn) {
		if err := handleConn(conn, service, syscall.Token(token), timeout); err != nil {
			logger.Error(err)
		}
	})
}

func doListen(server net.Listener, handler func(net.Conn)) error {
	for {
		conn, err := server.Accept()
		if err != nil {
			return err
		}
		logger.Info("New connection from", conn.RemoteAddr().String())
		go handler(conn)
	}
}

func handleConn(conn net.Conn, service string, token syscall.Token, timeout int) error {
	defer conn.Close()

	cmd := exec.Command(service)
	cmd.SysProcAttr = &syscall.SysProcAttr{Token: token}
	cmd.Stdout = conn

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	job, err := windows.NewJobFromProcess(cmd.Process)
	if err != nil {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return err
	}

	go func() {
		var buf [1]byte
		duration := time.Duration(timeout)
		for {
			conn.SetReadDeadline(time.Now().Add(duration * time.Second))
			if _, err := conn.Read(buf[:]); err != nil {
				logger.Warning(err)
				break
			}
			if _, err := stdin.Write(buf[:]); err != nil {
				logger.Warning(err)
				break
			}
		}

		if err := windows.TerminateJob(job); err != nil {
			logger.Error(err)
		}

		if cmd.Process != nil && (cmd.ProcessState != nil && !cmd.ProcessState.Exited()) {
			if err := cmd.Process.Kill(); err != nil {
				logger.Info(err)
			}
		}
	}()

	logger.Info(fmt.Sprintf("Started service \"%s\" with pid %d", service, cmd.Process.Pid))

	if err := cmd.Wait(); err != nil {
		logger.Warning(err)
	}

	return nil
}
