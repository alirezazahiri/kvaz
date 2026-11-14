package main

import (
	"fmt"
	"kvaz/internal/dto"
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
				return []byte(err.Error())
			}
			res := processor.Process(cmd)
			switch res.Kind {
			case dto.KindSimpleString, dto.KindBulkString:
				return []byte(res.Text)
			case dto.KindInteger:
				return fmt.Appendf([]byte{}, "%d", res.Integer)
			case dto.KindNull:
				return []byte("NULL")
			case dto.KindError:
				return []byte(res.Text)
			default:
				return []byte("NOOP")
			}
		})
	}
}
