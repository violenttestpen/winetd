package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

const integrityUsage = "Run service with assigned integrity level: %v"

var (
	bind      = "0.0.0.0"
	port      = 8080
	integrity = "Untrusted"
	verbosity = 0

	service string
)

func init() {
	i, integrityLevels := 0, make([]string, len(sidWinIntegrityLevels))
	for k := range sidWinIntegrityLevels {
		integrityLevels[i] = k
		i++
	}

	flag.StringVar(&bind, "bind", bind, "Address to bind to")
	flag.IntVar(&port, "port", port, "Port number to bind to")
	flag.StringVar(&integrity, "integrity", integrity, fmt.Sprintf(integrityUsage, integrityLevels))
	flag.IntVar(&verbosity, "verbosity", verbosity, "Verbosity mode (0-2)")
	flag.StringVar(&service, "server", service, "Path to service to be daemonized")
	flag.Parse()
}

func main() {
	if _, err := os.Stat(service); err != nil {
		log.Fatal(err)
	}

	sid, ok := sidWinIntegrityLevels[integrity]
	if !ok {
		log.Fatal("Invalid Integrity Level")
	}

	if err := beginListener(fmt.Sprintf("%s:%d", bind, port), service, sid, verbosity); err != nil {
		log.Fatal(err)
	}
}

func beginListener(bind, service string, sid uint32, verbosity int) error {
	server, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			return err
		}
		if verbosity >= 1 {
			log.Printf("New connection from %v\n", conn.RemoteAddr())
		}
		go func() {
			if err := handleConn(service, sid, conn, verbosity); err != nil {
				log.Println(err)
			}
		}()
	}
}

func handleConn(service string, sid uint32, conn net.Conn, verbosity int) error {
	defer conn.Close()

	token, err := getIntegrityLevelToken(sid)
	if err != nil {
		return err
	}

	cmd := exec.Command(service)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Token:         syscall.Token(token),
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | windows.CREATE_SUSPENDED,
	}
	cmd.Stdout = conn

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); verbosity >= 1 && err != nil {
		return err
	}
	token.Close()

	job, err := NewJobFromProcess(cmd.Process)
	if err != nil {
		cmd.Process.Kill()
		return err
	}
	if err := ResumeThread(uint32(cmd.Process.Pid)); err != nil {
		cmd.Process.Kill()
		return err
	}

	go func(w io.Writer, r io.Reader) {
		if _, err := io.Copy(w, r); verbosity >= 2 && err != nil {
			log.Println(err)
		}

		if err := TerminateJob(job); err != nil {
			log.Println(err)
		}

		if cmd.Process != nil && !cmd.ProcessState.Exited() {
			if err := cmd.Process.Kill(); verbosity >= 2 && err != nil {
				log.Println(err)
			}
		}
	}(stdin, conn)

	if verbosity >= 1 {
		log.Printf("Started service \"%s\" with pid %d\n", service, cmd.Process.Pid)
	}

	if err := cmd.Wait(); verbosity >= 1 && err != nil {
		log.Println(err)
	}

	return nil
}
