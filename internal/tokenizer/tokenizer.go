// Package tokenizer: this one is responsible for tokenization
package tokenizer

import (
	"strings"
	"unicode"
)

type Tokenizer struct{}

func New() *Tokenizer {
	return &Tokenizer{}
}

func (t *Tokenizer) SplitFirstToken(b []byte) (TokenType, string) {
	sb := strings.Builder{}
	clean := strings.TrimSpace(string(b))

	for i, r := range clean {
		if unicode.IsLetter(r) {
			sb.WriteRune(r)
		} else {
			clean = clean[i:]
			break
		}
	}

	tok := sb.String()
	tt, ok := keywords[tok]
	result := strings.TrimSpace(clean)
	if !ok {
		return ERROR, tok
	}

	return tt, result
}

func Tokenize(s string) []string {
	return strings.Split(s, " ")
}
