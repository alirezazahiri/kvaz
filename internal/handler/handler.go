// Package handler: this one is responsible for reading byte stream from the tcp connection
package handler

import (
	"errors"
	"fmt"
	"io"
	"kvaz/internal/tcp"
)

type HandlerFn func([]byte) []byte
type ConnectionHandler struct {
	tc *tcp.TCPConnection
}

func NewConnectionHandler(tc *tcp.TCPConnection) *ConnectionHandler {
	return &ConnectionHandler{tc: tc}
}

func (ch *ConnectionHandler) Connect(handler HandlerFn, options ...int) {
	bufsize := 1024
	defer ch.tc.Close()
	if len(options) > 0 {
		bufsize = options[0]
	}
	for {
		d := make([]byte, bufsize)
		n, err := ch.tc.Read(d)

		if err != nil {
			if errors.Is(err, io.EOF) {
				if err := ch.tc.CloseWrite(); err != nil {
					fmt.Println("connection closed with client")
				}
				return
			}
			fmt.Println("error occured while reading data:", err.Error())
			continue
		}

		bytes := d[:n]

		if n == 0 {
			return
		}

		response := handler(bytes)

		if _, err := ch.tc.Write(response); err != nil {
			fmt.Println("error occured while writing data")
			return
		}
	}
}
