package tcp

import (
	"fmt"
	"net"
)

type TCPServer struct {
	Port     uint16
	listener net.Listener
}

func NewTCPServer(port uint16) *TCPServer {
	l, err := listen(port)

	if err != nil {
		panic(err.Error())
	}

	return &TCPServer{
		Port:     port,
		listener: l,
	}
}

func listen(port uint16) (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))

	if err != nil {
		return nil, err
	}

	fmt.Printf("server is listening on port %d\n", port)

	return listener, nil
}

func (ts *TCPServer) Accept() (*TCPConnection, error) {
	conn, err := ts.listener.Accept()

	if err != nil {
		return nil, fmt.Errorf("failed to establish a connection: %s", err.Error())
	}

	return NewTCPConnection(conn), nil
}

func (ts *TCPServer) Close() error {
	return ts.listener.Close()
}
