package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running at localhost8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			fmt.Printf("Accept %v\n", conn.RemoteAddr())

			_, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			// _, err := http.ReadRequest(bufio.NewReaderSize(conn, 1))
			// if err != nil {
			// 	panic(err)
			// }
			// dump, err := httputil.DumpRequest(req, true)
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Println(string(dump))

			// resp := http.Response{
			// 	StatusCode: 200,
			// 	ProtoMajor: 1,
			// 	ProtoMinor: 0,
			// 	Body:       ioutil.NopCloser(strings.NewReader("Hello World\n")),
			// }
			// rdump, err := httputil.DumpResponse(&resp, true)
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Println(string(rdump))
			// resp.Write(conn)
			conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\nHello World\n"))
			conn.Close()
		}()
	}
}
