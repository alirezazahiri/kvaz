// Package dto: this one is responsible for the construction of data-transfer-object
package dto

type Kind uint8

const (
	KindSimpleString Kind = 1
	KindBulkString   Kind = 2
	KindInteger      Kind = 3
	KindNull         Kind = 4
	KindError        Kind = 5
)

type Response struct {
	Kind    Kind
	Text    string
	Integer int64
}

func SimpleString(s string) Response {
	return Response{Kind: KindSimpleString, Text: s}
}
func BulkString(s string) Response {
	return Response{Kind: KindBulkString, Text: s}
}
func Integer(i int64) Response {
	return Response{Kind: KindInteger, Integer: i}
}
func Null() Response {
	return Response{Kind: KindNull}
}
func Error(m string) Response {
	return Response{Kind: KindError, Text: m}
}
