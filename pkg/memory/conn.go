package memory

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func HandleConnection(listener net.Listener) {
	data := NewData()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection %v", err)
			continue
		}

		fmt.Printf("new tcp connection %s \n", conn.RemoteAddr())

		for {
			// conn.Write([]byte("Enter command: "))
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("connection closed")
					conn.Close()
					continue
				}
			}

			result, error := parseCommand(string(buf[:n]), data)
			if error != nil {
				conn.Write([]byte(error.Error()))
				continue
			}

			conn.Write([]byte(result))
		}
	}
}
