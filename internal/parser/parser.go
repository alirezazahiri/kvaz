// Package parser: this one is responsible for parsing commands
package parser

import (
	"fmt"
	"kvaz/internal/command"
	"kvaz/internal/tokenizer"
	"strconv"
)

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(tt tokenizer.TokenType, rest string) (command.Command, error) {
	args := tokenizer.Tokenize(rest)
	n := len(args)

	switch tt {
	case tokenizer.PING:
		if n > 0 {
			return nil, p.errorAt("PING")
		}
		return p.ping()
	case tokenizer.GET:
		if n > 1 {
			return nil, p.errorAt("GET", 1, n)
		}
		return p.get(args[0])
	case tokenizer.SET:
		if n != 2 {
			return nil, p.errorAt("SET", 2, n)
		}
		return p.set(args[0], args[1])
	case tokenizer.DEL:
		if n > 1 {
			return nil, p.errorAt("DEL", 1, n)
		}
		return p.del(args[0])
	case tokenizer.EXISTS:
		if n > 1 {
			return nil, p.errorAt("EXISTS", 1, n)
		}
		return p.exists(args[0])
	case tokenizer.EXPIRE:
		if n != 2 {
			return nil, p.errorAt("EXPIRE", 2, n)
		}
		secs, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return nil, p.errorAt("EXPIRE")
		} else if secs < 0 {
			return nil, fmt.Errorf("expected non-negative numbers, received: %d", secs)
		}
		return p.expire(args[0], secs)
	case tokenizer.TTL:
		if n > 1 {
			return nil, p.errorAt("TTL", 1, n)
		}
		return p.ttl(args[0])
	case tokenizer.ERROR:
		return nil, fmt.Errorf("unexpected token '%s'", rest)
	default:
		return nil, fmt.Errorf("unexpected parse error")
	}
}

func (p *Parser) errorAt(t string, rest ...int) error {
	n := len(rest)
	if n > 0 {
		if n != 2 {
			return fmt.Errorf("INTERNAL: invalid arguments")
		} else {
			return fmt.Errorf("expected %d after %s, but got %d", rest[0], t, rest[1])
		}
	}
	return fmt.Errorf("unexpected args for token %s", t)
}
