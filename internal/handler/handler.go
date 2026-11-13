package handler

import (
	"fmt"
	"kvaz/internal/tcp"
)

type ConnectionHandler struct {
	tc *tcp.TCPConnection
}

func NewConnectionHandler(tc *tcp.TCPConnection) *ConnectionHandler {
	return &ConnectionHandler{tc: tc}
}

func (ch *ConnectionHandler) Connect() {
	for {
		d := make([]byte, 1024)
		if _, err := ch.tc.Read(d); err != nil {
			fmt.Println("error occured while reading data")
			continue
		}
		fmt.Print("data received: ", string(d))

		if _, err := ch.tc.Write(d); err != nil {
			fmt.Println("error occured while writing data")
			continue
		}
	}
}
