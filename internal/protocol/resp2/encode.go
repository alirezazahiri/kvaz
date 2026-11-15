package resp2

import "strconv"

var crlf = []byte("\r\n")

func Encode(v Value) ([]byte, error) {
	return New().encode(v)
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
		return nil, ProtocolErrorUnknownValueKind
	}
}
