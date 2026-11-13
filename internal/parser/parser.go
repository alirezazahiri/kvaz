// Package parser: this one is responsible for parsing commands
package parser

import (
	"fmt"
	"kvaz/internal/command"
	"kvaz/internal/tokenizer"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(tt tokenizer.TokenType, rest string) (command.Command, error) {
	switch tt {
	case tokenizer.PING:
		return p.ping()
	case tokenizer.GET:
		return p.get()
	case tokenizer.SET:
		return p.set()
	case tokenizer.DEL:
		return p.del()
	case tokenizer.EXISTS:
		return p.exists()
	case tokenizer.EXPIRE:
		return p.expire()
	case tokenizer.TTL:
		return p.ttl()
	case tokenizer.ERROR:
		return nil, fmt.Errorf("unexpected token '%s'", rest)
	default:
		return nil, fmt.Errorf("unexpected parse error")
	}
}
