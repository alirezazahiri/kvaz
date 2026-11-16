package resp2

type Kind uint8

const (
	KindSimpleString Kind = 1
	KindBulkString   Kind = 2
	KindInteger      Kind = 3
	KindNull         Kind = 4
	KindError        Kind = 5
	KindArray        Kind = 6
)

type Value struct {
	Kind     Kind
	Text     string
	Integer  int64
	Bytes    []byte
	Elements []Value
}

func SimpleString(s string) Value {
	return Value{Kind: KindSimpleString, Text: s}
}
func BulkString(s string) Value {
	return Value{Kind: KindBulkString, Bytes: []byte(s)}
}
func Integer(i int64) Value {
	return Value{Kind: KindInteger, Integer: i}
}
func Null() Value {
	return Value{Kind: KindNull}
}
func Error(m string) Value {
	return Value{Kind: KindError, Text: m}
}
func Array(v []Value) Value {
	return Value{Kind: KindArray, Elements: v, Integer: int64(len(v))}
}
