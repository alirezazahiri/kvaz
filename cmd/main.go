package main

import (
	"fmt"
	"kvaz/internal/command"
	"kvaz/internal/handler"
	"kvaz/internal/parser"
	"kvaz/internal/tcp"
	"kvaz/internal/tokenizer"
)

func main() {
	ts := tcp.NewTCPServer(6379)
	defer ts.Close()
	cmd := make(chan command.Command, 10)
	go func() {
		for {
			task := <-cmd
			fmt.Printf("executing task: %#v\n", task)
			task.Execute()
		}
	}()

	for {
		tc, err := ts.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		ch := handler.NewConnectionHandler(tc)
		t := tokenizer.New()
		p := parser.New()
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
