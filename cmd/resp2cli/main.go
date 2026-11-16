package main

import (
	"fmt"

	resp2handler "kvaz/internal/handler/resp2"
	resp2parser "kvaz/internal/parser/resp2"
	"kvaz/internal/processor"
	"kvaz/internal/storage"
	"kvaz/internal/tcp"
)

func main() {
	ts := tcp.NewTCPServer(6379)
	defer ts.Close()

	store := storage.NewInMemoryStorage()
	proc := processor.New(store)
	parser := resp2parser.New()

	for {
		conn, err := ts.Accept()
		if err != nil {
			fmt.Println("accept:", err)
			continue
		}

		handler := resp2handler.New(conn, parser, proc)

		go func() {
			if err := handler.Serve(); err != nil {
				fmt.Println("connection:", err)
			}
		}()
	}
}
