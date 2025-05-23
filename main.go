package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"syscall"
)

const BUFFER_SIZE = 4096

func read_bytes(buf *[]byte, fd int, sb *strings.Builder) int {
	r := 0
	for ok := true; ok; {
		n, err := syscall.Read(fd, *buf)
		if err != nil {
			log.Fatal(err)
		}

		r += n
		(*sb).WriteString(string((*buf)[:n]))
		ok = (n != 0)
	}
	return r
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, BUFFER_SIZE)
	headers := make([]byte, BUFFER_SIZE)

	for {
		n, err := conn.Read(headers)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		// fmt.Printf("%s\n", string(headers[:n]))
		headers := strings.Split(string(headers[:n]), "\r\n")

		parsed_header := strings.Split(headers[0], " ")

		req_type := parsed_header[0]
		infile := parsed_header[1][1:]

		fmt.Println(req_type, infile)
		fmt.Println(len(infile))

		fd, err := syscall.Open("index.html", syscall.O_RDONLY, 0)

		status_code := 200

		if err != nil {
			status_code = 404
		}
		defer syscall.Close(fd)

		var sb strings.Builder
		bytes := read_bytes(&buffer, fd, &sb)

		response := fmt.Sprintf("HTTP/1.1 %d OK\r\n"+
			"Content-Type: text/plain\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n%s", status_code, bytes, sb.String())

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Server listening on port 8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		handleConnection(conn)
	}
}
