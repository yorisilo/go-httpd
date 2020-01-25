package main

import (
	"log"
	"net"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8080")
	if err != nil {
		panic(err)
	}
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		log.Printf("accept\n")

		go func() {
			request := make([]byte, 78)
			conn.Read(request)
			// http.ReadRequest(bufio.NewReader(conn))
			conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\nHello World\n"))
			conn.Close()
		}()
	}
}
