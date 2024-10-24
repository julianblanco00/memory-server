package memory

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func handleReadFromConn(conn net.Conn, data *Data) {
	for {
		// conn.Write([]byte("Enter command: "))
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("connection closed")
				conn.Close()
				break
			}
		}

		connId := buf[:36]
		cmd := []byte{}

		for _, v := range buf[37:] {
			if v == 0 {
				break
			}
			cmd = append(cmd, v)
		}

		result, error := parseCommand(string(cmd), data)
		if error != nil {
			conn.Write([]byte(error.Error()))
			continue
		}

		fmt.Fprint(conn, string(connId), result)
	}
}

func HandleConnection(listener net.Listener) {
	data := NewData()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection %v", err)
			continue
		}

		fmt.Printf("new tcp connection %s \n", conn.RemoteAddr())

		go handleReadFromConn(conn, data)
	}
}
