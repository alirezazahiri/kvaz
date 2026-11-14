package main

import (
	"fmt"
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

	for {
		tc, err := ts.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		ch := handler.NewConnectionHandler(tc)

		go ch.Connect(func(b []byte) []byte {
			cmd, err := p.Parse(t.SplitFirstToken(b))
			if err != nil {
				return []byte(err.Error() + "\n")
			}
			return []byte(processor.Process(cmd) + "\n")
		})
	}
}
