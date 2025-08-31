package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var CRLF = "\r\n"

func main() {
	fmt.Println("Logs from HTTP Server")

	l, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	path, err := extractURLPath(conn)
	if err != nil {
		fmt.Println("Error extracting url path from request: ", err)
		os.Exit(1)
	}

	resp := handlePath(path)

	conn.Write(resp)
}

func extractURLPath(conn net.Conn) (string, error) {
	rcvBuff := make([]byte, 1024)

	_, err := conn.Read(rcvBuff)

	if err != nil {
		return "", err
	}

	req := string(rcvBuff)

	lines := strings.Split(req, CRLF)
	path := strings.Split(lines[0], " ")[1]

	return path, nil
}

func handlePath(path string) []byte {
	if path == "/" {
		return []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if str, ok := strings.CutPrefix(path, "/echo/"); ok {
		resp := "HTTP/1.1 200 OK" + CRLF
		resp += "Content-Type: text/plain" + CRLF
		resp += fmt.Sprintf("Content-Length: %v", len(str)) + CRLF
		resp += CRLF
		resp += str
		return []byte(resp)
	} else {
		return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
}
