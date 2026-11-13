package tokenizer

var keywords = map[string]TokenType{
	"PING":   PING,
	"SET":    SET,
	"GET":    GET,
	"DEL":    DEL,
	"EXISTS": EXISTS,
	"EXPIRE": EXPIRE,
	"TTL":    TTL,
}
