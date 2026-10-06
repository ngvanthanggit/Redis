package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println(conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	requestCount := 0
	for {
		request, err := http.ReadRequest(reader)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return // Client disconnected or send invalid request
		}
		requestCount += 1

		io.Copy(io.Discard, request.Body)
		request.Body.Close()

		time.Sleep(time.Second * 5)

		body := fmt.Sprintf("Hello World %d\r\n", requestCount)
		connection := "keep-alive"

		if request.Close {
			connection = "close"
		}

		response := fmt.Sprintf(
			"HTTP/1.1 200 OK\r\n"+
				"Content-Length: %d\r\n"+
				"Content-Type: text/plain\r\n"+
				"Connection: %s\r\n"+
				"\r\n"+
				"%s",
			len(body),
			connection,
			body,
		)

		if _, err := conn.Write([]byte(response)); err != nil {
			return // Cliend disconnected while receiving the response
		}

		if request.Close {
			return // Client requested Connection: close
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	for {
		// conn == socket
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// create a new goroutine to handle multiple connections
		// multiple clients can send request at the same time
		go handleConnection(conn)
	}
}
