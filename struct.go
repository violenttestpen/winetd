package main

// Service represents an internet daemon instance
type Service struct {
	Bind     string `json:"bind"`
	Protocol string `json:"protocol"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"`
	Disable  bool   `json:"disable"`
	Wait     bool   `json:"wait"`
	// User       string `json:"user"`
}

var defaultService = Service{
	Bind:     "127.0.0.1",
	Protocol: "tcp",
	Server:   "cmd.exe",
	Port:     8080,
	Timeout:  60,
	Disable:  true,
	Wait:     false,
}
