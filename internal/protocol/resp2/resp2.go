package resp2

import (
	"errors"
	"strconv"
)

var crlf = []byte("\r\n")

var errUnknownValueKind = errors.New("unknown value kind")

type ProtocolError string

const (
	ProtocolErrorInsufficientInputBytes ProtocolError = "insufficient input bytes"
	ProtocolErrorInputMalformed         ProtocolError = "input is malformed"
)

func (e ProtocolError) Error() string {
	return string(e)
}

type Protocol struct {
	buffer []byte
	cursor int
}

func New() *Protocol {
	return &Protocol{}
}

func (p *Protocol) Feed(b []byte) {
	p.buffer = append(p.buffer, b...)
}

func (p *Protocol) Next() (value Value, success bool, pErr error) {
	return p.next(true)
}

func (p *Protocol) next(compactOnSuccess bool) (value Value, success bool, pErr error) {
	value, success, pErr = Value{}, true, nil
	if p.cursor >= len(p.buffer) {
		success = false
		pErr = ProtocolErrorInsufficientInputBytes
		return
	}
	switch p.buffer[p.cursor] {
	case '+':
		b, err := p.scan()
		if err != nil {
			success = false
			pErr = err
			break
		}
		value = SimpleString(string(b))
	case '-':
		b, err := p.scan()
		if err != nil {
			success = false
			pErr = err
			break
		}
		value = Error(string(b))
	case ':':
		b, err := p.scan()
		if err != nil {
			success = false
			pErr = err
			break
		}
		num, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			success = false
			pErr = ProtocolErrorInputMalformed
			break
		}
		value = Integer(num)
	case '$':
		b, err := p.scan()
		if err != nil {
			success = false
			pErr = err
			break
		}
		length, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			success = false
			pErr = ProtocolErrorInputMalformed
			break
		}
		if length == -1 {
			value = Null()
			break
		}
		if length < -1 {
			success = false
			pErr = ProtocolErrorInputMalformed
			break
		}
		payload, err := p.scanBulkPayload(int(length))
		if err != nil {
			success = false
			pErr = err
			break
		}
		value = BulkString(string(payload))
	case '*':
		b, err := p.scan()
		if err != nil {
			success = false
			pErr = err
			break
		}
		size, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			success = false
			pErr = ProtocolErrorInputMalformed
			break
		}
		if size == -1 {
			value = Null()
			break
		}
		if size < -1 {
			success = false
			pErr = ProtocolErrorInputMalformed
			break
		}
		values := make([]Value, size)
		for i := range size {
			value, success, pErr = p.next(false)
			if !success || pErr != nil {
				break
			}
			values[i] = value
		}
		if success {
			value = Array(values)
		}
	default:
		success = false
		pErr = ProtocolErrorInputMalformed
	}
	if success {
		if compactOnSuccess {
			p.compactBuffer()
		}
	} else {
		p.cursor = 0
	}
	return
}
func (p *Protocol) Reset() {
	p.buffer = make([]byte, 0)
	p.cursor = 0
}
func (p *Protocol) compactBuffer() {
	n := len(p.buffer)
	if p.cursor >= n {
		p.Reset()
		return
	}
	p.buffer = p.buffer[p.cursor:]
	p.cursor = 0
}

func (p *Protocol) scanBulkPayload(n int) ([]byte, error) {
	end := p.cursor + n
	if end > len(p.buffer) {
		return nil, ProtocolErrorInsufficientInputBytes
	}
	payload := p.buffer[p.cursor:end]
	p.cursor = end
	if p.cursor+1 >= len(p.buffer) {
		return nil, ProtocolErrorInsufficientInputBytes
	}
	if p.buffer[p.cursor] != '\r' || p.buffer[p.cursor+1] != '\n' {
		return nil, ProtocolErrorInputMalformed
	}
	p.cursor += 2
	return payload, nil
}

// readUntilNextCRLF reads until it sees \r\n and leaves the cursor past \n.
func (p *Protocol) readUntilNextCRLF() bool {
	for p.cursor < len(p.buffer) {
		if p.buffer[p.cursor] != '\r' {
			p.cursor++
			continue
		}
		if p.cursor+1 >= len(p.buffer) {
			return false
		}
		if p.buffer[p.cursor+1] != '\n' {
			return false
		}
		p.cursor += 2
		return true
	}
	return false
}

func (p *Protocol) scan() ([]byte, error) {
	start := p.cursor + 1
	if !p.readUntilNextCRLF() {
		return nil, ProtocolErrorInsufficientInputBytes
	}
	return p.buffer[start : p.cursor-2], nil
}

// Encode serializes a RESP2 value.
func Encode(v Value) ([]byte, error) {
	return (&Protocol{}).encode(v)
}

func (p *Protocol) encode(v Value) ([]byte, error) {
	switch v.Kind {
	case KindSimpleString:
		b := append([]byte{'+'}, v.Text...)
		return append(b, crlf...), nil
	case KindError:
		b := append([]byte{'-'}, v.Text...)
		return append(b, crlf...), nil
	case KindInteger:
		b := append([]byte{':'}, strconv.AppendInt(nil, v.Integer, 10)...)
		return append(b, crlf...), nil
	case KindNull:
		return []byte("$-1\r\n"), nil
	case KindBulkString:
		payload := v.Bytes
		if payload == nil {
			payload = []byte(v.Text)
		}
		b := append([]byte{'$'}, strconv.AppendInt(nil, int64(len(payload)), 10)...)
		b = append(b, crlf...)
		b = append(b, payload...)
		return append(b, crlf...), nil
	case KindArray:
		b := append([]byte{'*'}, strconv.AppendInt(nil, int64(len(v.Elements)), 10)...)
		b = append(b, crlf...)
		for _, el := range v.Elements {
			part, err := p.encode(el)
			if err != nil {
				return nil, err
			}
			b = append(b, part...)
		}
		return b, nil
	default:
		return nil, errUnknownValueKind
	}
}