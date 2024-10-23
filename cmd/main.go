package main

import (
	"custom-redis/pkg/memory"
	"fmt"
	"log"
	"net"
	"time"
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
	go testClient()
	initServer("127.0.0.1", "4444")
}

func testClient() {
	time.Sleep(1 * time.Second)
	fmt.Println("connecting client...")
	client, err := net.Dial("tcp", ":4444")
	if err != nil {
		log.Fatal("error connecting client")
	}

	go func() {
		for {
			buf := make([]byte, 1024)
			client.Read(buf)
			fmt.Println("response from server: ", string(buf))
		}
	}()

	fmt.Fprintf(
		client,
		"SET key1 val2 NX GET KEEPTTL",
		time.Now().Add(time.Minute).UnixMilli(),
	)
	// time.Sleep(1 * time.Second)
	// client.Write([]byte("GET key1"))
	// time.Sleep(1 * time.Second)
	// client.Write([]byte("DEL key1 key2 key3"))
}
