package main

import (
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

const (
	BIND    = "127.0.0.1:8080"
	SERVICE = "./servers/echo.exe"
)

var listenerIsLoaded = false

func simulateListener() {
	if !listenerIsLoaded {
		go beginListener(BIND, SERVICE, sidWinIntegrityLevels["Untrusted"], 0)
		conn, err := net.Dial("tcp", BIND)
		if err != nil {
			log.Fatal("Unable to connect to listener")
		}
		conn.Close()
	}
	listenerIsLoaded = true
}

func simulateConnection(sleep int, wg *sync.WaitGroup) error {
	defer wg.Done()
	conn, err := net.Dial("tcp", BIND)
	if err != nil {
		return err
	}
	defer conn.Close()
	time.Sleep(time.Duration(sleep) * time.Second)
	return nil
}

func simulateConnections(n int, sleep int) {
	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go simulateConnection(sleep, wg)
	}
	wg.Wait()
}

func BenchmarkService10(b *testing.B) {
	simulateListener()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		simulateConnections(10, 10)
	}
}

func BenchmarkService100(b *testing.B) {
	simulateListener()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		simulateConnections(100, 10)
	}
}
