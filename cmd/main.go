package main

import (
	"online-store/pkg/httpserver"
	"time"
)

func main() {
	serverDeps := &httpserver.ServerDeps{
		Host:    "localhost",
		Port:    ":8080",
		Timeout: 5 * time.Second,
	}

	server := httpserver.NewServer(serverDeps)
	server.Start()
}
