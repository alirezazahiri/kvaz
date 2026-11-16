// Package resp2handler: this one is responsible for resp2 connections
package resp2handler

import (
	"errors"
	"fmt"
	"io"

	"kvaz/internal/dto"
	resp2parser "kvaz/internal/parser/resp2"
	"kvaz/internal/processor"
	protocol "kvaz/internal/protocol/resp2"
	"kvaz/internal/tcp"
)

type ConnectionHandler struct {
	conn      *tcp.TCPConnection
	protocol  *protocol.Protocol
	parser    *resp2parser.Parser
	processor *processor.Processor
}

func New(
	conn *tcp.TCPConnection,
	parser *resp2parser.Parser,
	processor *processor.Processor,
) *ConnectionHandler {
	return &ConnectionHandler{
		conn:      conn,
		protocol:  protocol.New(),
		parser:    parser,
		processor: processor,
	}
}

// Serve handles the complete lifecycle of one client connection.
func (h *ConnectionHandler) Serve() error {
	defer h.conn.Close()

	buffer := make([]byte, 4096)

	for {
		n, readErr := h.conn.Read(buffer)

		if n > 0 {
			h.protocol.Feed(buffer[:n])

			if err := h.processBuffered(); err != nil {
				return err
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return nil
			}
			return fmt.Errorf("read connection: %w", readErr)
		}

		if n == 0 {
			return nil
		}
	}
}

// processBuffered processes every complete RESP value currently buffered.
func (h *ConnectionHandler) processBuffered() error {
	for {
		value, success, err := h.protocol.Next()

		if err != nil {
			if errors.Is(err, protocol.ProtocolErrorInsufficientInputBytes) {
				// Wait for the next TCP read.
				return nil
			}

			_ = h.writeValue(protocol.Error(err.Error()))
			return fmt.Errorf("decode RESP2: %w", err)
		}

		if !success {
			return fmt.Errorf("RESP2 decoder returned no value without an error")
		}

		commands, err := h.parser.Parse(value)
		if err != nil {
			if writeErr := h.writeValue(protocol.Error(err.Error())); writeErr != nil {
				return writeErr
			}
			continue
		}

		for _, cmd := range commands {
			response := h.processor.Process(cmd)

			if err := h.writeValue(toRESP2(response)); err != nil {
				return err
			}
		}
	}
}

func toRESP2(response dto.Response) protocol.Value {
	switch response.Kind {
	case dto.KindSimpleString:
		return protocol.SimpleString(response.Text)

	case dto.KindBulkString:
		return protocol.BulkString(response.Text)

	case dto.KindInteger:
		return protocol.Integer(response.Integer)

	case dto.KindNull:
		return protocol.Null()

	case dto.KindError:
		return protocol.Error(response.Text)

	default:
		return protocol.Error("unsupported response type")
	}
}

func (h *ConnectionHandler) writeValue(value protocol.Value) error {
	data, err := h.protocol.Encode(value)
	if err != nil {
		return fmt.Errorf("encode RESP2 response: %w", err)
	}

	if err := writeAll(h.conn, data); err != nil {
		return fmt.Errorf("write RESP2 response: %w", err)
	}

	return nil
}

// writeAll handles short writes. A successful Write isn't guaranteed
// to write the entire supplied byte slice.
func writeAll(conn *tcp.TCPConnection, data []byte) error {
	for len(data) > 0 {
		n, err := conn.Write(data)

		if err != nil {
			return err
		}

		if n == 0 {
			return io.ErrShortWrite
		}

		data = data[n:]
	}

	return nil
}
