package main

import (
	"fmt"
	"kvaz/internal/command"
	"kvaz/internal/handler"
	"kvaz/internal/parser"
	"kvaz/internal/processor"
	"kvaz/internal/storage"
	"kvaz/internal/tcp"
	"kvaz/internal/tokenizer"
)

const PoolSize = 10

func main() {
	ts := tcp.NewTCPServer(6379)
	defer ts.Close()

	storage := storage.NewInMemoryStorage()
	processor := processor.New(storage)

	t := tokenizer.New()
	p := parser.New()

	cmd := make(chan command.Command, PoolSize)

	go func() {
		for {
			processor.Process(<-cmd)
		}
	}()

	for {
		tc, err := ts.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		ch := handler.NewConnectionHandler(tc)

		go ch.Connect(func(b []byte) {
			v, err := p.Parse(t.SplitFirstToken(b))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			cmd <- v
		})
	}
}
