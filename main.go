package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
)

const BUFFER_SIZE = 4096

func read_bytes(buf *[]byte, infile string) {
	fmt.Printf("Reading %s\n", infile)
	fd, err := syscall.Open(infile, syscall.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer syscall.Close(fd)

	for ok := true; ok; {
		n, err := syscall.Read(fd, *buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", string((*buf)[:n]))
		ok = (n != 0)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, BUFFER_SIZE)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		read_bytes(&buffer, "index.html")
		fmt.Printf("%s\n", string(buffer[:n]))
		body := "Hello from TCP server"
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n%s", len(body), body)
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
