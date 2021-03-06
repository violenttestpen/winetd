package main

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	BIND      = "127.0.0.1:8080"
	SERVICE   = "./servers/echoserver.exe"
	INTEGRITY = "Untrusted"
)

var onceListener sync.Once

func simulateListener() {
	svc := Service{Bind: BIND, Server: SERVICE, Timeout: 5}
	go run(svc, INTEGRITY)
}

func simulateConnection(timeout int) error {
	conn, err := net.Dial("tcp", BIND)
	if err != nil {
		return err
	}
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	io.WriteString(conn, "hello world\r\n")
	io.ReadAll(conn)
	return conn.Close()
}

func simulateConnections(n int, timeout int) {
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			simulateConnection(timeout)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkService1(b *testing.B) {
	onceListener.Do(simulateListener)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		simulateConnection(1)
	}
}

func BenchmarkService10(b *testing.B) {
	onceListener.Do(simulateListener)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		simulateConnections(10, 5)
	}
}

func TestListener(t *testing.T) {
	var err error
	assert := assert.New(t)

	invalidSvc := Service{Bind: BIND, Server: "invalid.service"}
	err = run(invalidSvc, INTEGRITY)
	assert.Error(err)

	nonExistentSvc := Service{Bind: BIND, Server: "NonExistent"}
	err = run(nonExistentSvc, SERVICE)
	assert.Error(err)

	go func() {
		svc := Service{Bind: BIND, Server: SERVICE, Timeout: 5}
		err = run(svc, INTEGRITY)
		assert.Error(err)
	}()
	time.Sleep(5 * time.Second)
	err = simulateConnection(5)
	assert.NoError(err)
}
