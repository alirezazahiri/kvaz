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
			res := processor.Process(cmd)
			if res.Error != nil {
				return []byte(res.Error.Error() + "\n")
			}
			if res.Value != nil {
				return fmt.Appendf([]byte{}, "value: %v, status: %s\n", res.Value, res.Status)
			}
			return fmt.Appendf([]byte{}, "status: %s\n", res.Status)
		})
	}
}
