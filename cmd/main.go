package main

import (
	"fmt"
	"kvaz/internal/handler"
	"kvaz/internal/tcp"
)

func main() {
	ts := tcp.NewTCPServer(6379)
	defer ts.Close()

	for {
		tc, err := ts.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		ch := handler.NewConnectionHandler(tc)
		go ch.Connect()
	}
}
