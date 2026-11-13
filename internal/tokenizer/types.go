package tokenizer

type TokenType uint8

const (
	PING TokenType = iota + 1
	SET
	GET
	DEL
	EXISTS
	EXPIRE
	TTL

	ERROR
)
