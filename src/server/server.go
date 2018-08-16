package server

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func Serve(host, port string, handler func(*http.Request) (http.Header, string)) {
	if port != "" {
		host = host + ":" + port
	}
	ln, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on " + host + ".")
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		request, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			panic(err)
		}
		header, body := handler(request)
		response := http.Response{
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 0,
			Body:       ioutil.NopCloser(strings.NewReader(body)),
		}
		response.Header = header
		response.Write(conn)
		conn.Close()
	}
}
