package memory

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const CONN_ID_LENGTH = 16

var (
	sData *StringData
	hData *HashData
)

func handleReadFromConn(conn net.Conn) {
	sData = NewStringData()
	hData = NewHashData()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("connection closed")
				conn.Close()
				break
			}
		}

		if n < 16 {
			fmt.Println("received an invalid buffer", n, err)
			conn.Close()
			break
		} else {
			fmt.Printf("new tcp connection %s \n", conn.RemoteAddr())
		}

		connId := buf[:CONN_ID_LENGTH]
		cmd := buf[CONN_ID_LENGTH:n]

		result, error := parseCommand(strings.TrimSpace(string(cmd)))
		if error != nil {
			fmt.Fprint(conn, string(connId), error)
			continue
		}

		if result == nil {
			conn.Write(connId)
		} else {
			var r []byte

			switch v := result.(type) {
			case int:
				r = []byte(strconv.Itoa(v))
			case string:
				r = []byte(v)
			default:
			}

			bytes := append(connId, r...)

			conn.Write(bytes)
		}
	}
}

func HandleConnection(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection %v", err)
			continue
		}

		go handleReadFromConn(conn)
	}
}
