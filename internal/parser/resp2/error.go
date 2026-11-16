package resp2parser

import "fmt"

type ParseError string

const (
	ParseErrorInvalid                    ParseError = "invalid"
	ParseErrorUnknownCommand             ParseError = "unknown command"
	ParseErrorEmptyArray                 ParseError = "array is empty"
	ParseErrorHeaderCannotContainInteger ParseError = "header cannot contain integer"
	ParseErrorArity                      ParseError = "wrong number of arguments for %s: expected %d arguments, but got %d"
	ParseErrorBadValue                   ParseError = "bad value received"
	ParseErrorInvalidBatchElement        ParseError = "invalid batch element at index %d: expected RESP array"
	ParseErrorInvalidCommandAtBatch      ParseError = "invalid command at batch index %d: %w"
	ParseErrorExpectedArray              ParseError = "expected RESP array, got kind %#v"
	ParseErrorInvalidArgument            ParseError = "invalid argument at index %d: %w"
	ParseErrorExpectedCommandArray       ParseError = "expected command array"
)

func (e ParseError) Error() string {
	return string(e)
}

func (e ParseError) format(args ...any) error {
	return fmt.Errorf(string(e), args...)
}
