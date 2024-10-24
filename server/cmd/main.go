package main

import (
	"custom-redis/pkg/memory"
	"fmt"
	"log"
	"net"
)

func initServer(addr, p string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", addr, p))
	if err != nil {
		log.Fatal("error starting server")
	}

	defer listener.Close()

	memory.HandleConnection(listener)
}

func main() {
	initServer("127.0.0.1", "4444")
}
