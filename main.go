package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"syscall"
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
	sid, ok := sidWinIntegrityLevels[integrity]
	if !ok {
		log.Fatal("Invalid Integrity Level")
	}

	if err := beginListener(fmt.Sprintf("%s:%d", bind, port), service, sid, verbosity); err != nil {
		log.Fatal(err)
	}
}

func beginListener(bind, service, sid string, verbosity int) error {
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
		go handleConn(service, sid, conn, verbosity)
	}
}

func handleConn(service, sid string, conn net.Conn, verbosity int) {
	defer conn.Close()

	token, err := getIntegrityLevelToken(sid)
	if err != nil {
		log.Println(err)
		return
	}
	defer token.Close()

	cmd := exec.Command(service)
	cmd.SysProcAttr = &syscall.SysProcAttr{Token: token}
	cmd.Stdout = conn

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		return
	}

	go func(w io.Writer, r io.Reader) {
		if _, err := io.Copy(w, r); verbosity >= 2 && err != nil {
			log.Println(err)
		}

		if cmd.Process != nil && !cmd.ProcessState.Exited() {
			if err := cmd.Process.Kill(); verbosity >= 2 && err != nil {
				log.Println(err)
			}
		}
	}(stdin, conn)

	if err := cmd.Start(); verbosity >= 1 && err != nil {
		log.Println(err)
		return
	}

	if verbosity >= 1 {
		log.Printf("Started service \"%s\" with pid %d\n", service, cmd.Process.Pid)
	}

	if err := cmd.Wait(); verbosity >= 1 && err != nil {
		log.Println(err)
	}
}
