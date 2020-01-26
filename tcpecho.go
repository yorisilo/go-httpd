package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	l, err := net.Listen("tcp4", "0.0.0.0:8080")
	defer l.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("wait...")
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			// 1024 byte までの文字列しか受け取れない仕様
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					return
				}
				panic(err)
			}
			_, err = conn.Write(buf[:n])
			if err != nil {
				panic(err)
			}
			conn.Close()
		}()
	}
}
