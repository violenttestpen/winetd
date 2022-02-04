package main

// Service represents an internet daemon instance
type Service struct {
	Bind       string `json:"bind"`
	Disable    bool   `json:"disable"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	Server     string `json:"server"`
	SocketType string `json:"socket_type"`
	Timeout    int    `json:"timeout"`
	// User       string `json:"user"`
	// Wait       bool   `json:"wait"`
}

var defaultService = Service{
	Bind:       "127.0.0.1",
	Disable:    true,
	Protocol:   "tcp",
	Port:       8080,
	Server:     "cmd.exe",
	SocketType: "stream",
	Timeout:    60,
}
