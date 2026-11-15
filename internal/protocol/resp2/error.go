package resp2

type ProtocolError string

const (
	ProtocolErrorInsufficientInputBytes ProtocolError = "insufficient input bytes"
	ProtocolErrorInputMalformed         ProtocolError = "input is malformed"
	ProtocolErrorUnknownValueKind       ProtocolError = "unknown value kind"
)

func (e ProtocolError) Error() string {
	return string(e)
}
