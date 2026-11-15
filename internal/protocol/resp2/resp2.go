// Package resp2: this on contains the implementations of resp2 protocol
package resp2

type Protocol struct {
	buffer []byte
	cursor int
}

func New() *Protocol {
	return &Protocol{}
}

