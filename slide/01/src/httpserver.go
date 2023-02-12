package main

import (
	"io"
	"net/http"
)

type Fuga struct {
	hoge int64
	fuga int64
	str  string
}

func handler(w http.ResponseWriter, r *http.Request) {
	// dump, err := httputil.DumpRequest(r, true)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(dump))
	// w.Write([]byte("HTTP/1.0 200 OK\r\n\r\nHello World\n"))
	io.WriteString(w, "Hello World\n")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
