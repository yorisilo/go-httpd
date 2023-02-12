package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"strings"
)

func main() {
	if os.Getenv("naive") == "on" {
		naiveHTTP()
	} else if os.Getenv("mini") == "on" {
		minimalHTTP()
	} else if os.Getenv("minic") == "on" {
		minimalClient()
	} else if os.Getenv("http10") == "on" {
		http10()
	}
}

func naiveHTTP() {
	ln, _ := net.Listen("tcp", "localhost:8888")
	fmt.Printf("Server is runnning at localhost:8888")
	conn, _ := ln.Accept()
	req := make([]byte, 1024)
	conn.Read(req)
	conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\nHello World\n"))
	conn.Close()
}

// X / _ /) < curl localhost:8888
// Hello World
// curl: (56) Recv failure: Connection reset by peer
// curl は http の仕様に則った通信を行うので、request を読まない場合、上記のようなエラーが出る
// conn.Close() は、 request メッセージを全部読んでなかった場合 RST パケットを client へ送る
func minimalHTTP() {
	ln, _ := net.Listen("tcp", "localhost:8888")
	conn, _ := ln.Accept()
	buf := make([]byte, 5)
	conn.Read(buf)
	conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\nHello World\n"))
	conn.Close()
}

func minimalClient() {
	conn, _ := net.Dial("tcp", "localhost:8888")
	conn.Write([]byte("GET / HTTP/1.0\r\n"))
	io.Copy(os.Stdout, conn)
	conn.Close()
}

func http10() {
	ln, _ := net.Listen("tcp", "localhost:8888")
	fmt.Printf("Server is runnning at localhost:8888\n")
	conn, _ := ln.Accept()
	buf := make([]byte, 5)

	// content-length の有無で分岐して、なければ、body は読まない。あればその文字数分だけ body を読むようにする
	reader := bufio.NewReader(conn)
	scanner := textproto.NewReader(reader)
	for {
		line, _ := scanner.ReadLine()

		if line == "" {
			break
		}
		hl := strings.Fields(line)
	}

	conn.Close()
}
