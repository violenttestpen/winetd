package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
)

func main() {
	bind := flag.String("bind", "0.0.0.0", "Address to bind to")
	port := flag.Int("port", 8080, "Port number to bind to")
	service := flag.String("server", "", "Path to service to be daemonized")
	username := flag.String("user", "", "The user to run the service as")
	password := flag.String("pass", "", "The password of the user")
	verbosity := flag.Int("verbosity", 0, "Verbosity mode (0-2)")
	flag.Parse()

	if err := beginListener(fmt.Sprintf("%s:%d", *bind, *port),
		*service, *username, *password, *verbosity); err != nil {
		log.Fatal(err)
	}
}

func beginListener(bind, service, username, password string, verbosity int) error {
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
		go handleConn(service, username, password, conn, verbosity)
	}
}

func handleConn(service, username, password string, conn net.Conn, verbosity int) {
	defer conn.Close()
	cmd := exec.Command(service)
	cmd.Stdout = conn

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		return
	}

	if len(username) > 0 {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		err := impersonate(username, password)
		if err != nil {
			log.Fatal(err)
		}
		defer mustRevertToSelf()
		if verbosity >= 2 {
			log.Printf("Impersonated as user %s\n", username)
		}
	}

	go func(w io.Writer, r io.Reader) {
		if _, err := io.Copy(w, r); verbosity >= 2 && err != nil {
			log.Println(err)
		}

		if cmd.Process != nil {
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
