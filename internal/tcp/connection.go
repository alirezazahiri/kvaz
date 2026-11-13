// this one is responsible for creating tcp connections
package tcp

import "net"

type TCPConnection struct {
	connection net.Conn
}

func NewTCPConnection(conn net.Conn) *TCPConnection {
	return &TCPConnection{
		connection: conn,
	}
}

func (tc *TCPConnection) Read(b []byte) (int, error) {
	return tc.connection.Read(b)
}

func (tc *TCPConnection) Write(b []byte) (int, error) {
	return tc.connection.Write(b)
}

func (tc *TCPConnection) Close() error {
	return tc.connection.Close()
}
