// Package resp2parser is responsible for parsing RESP2 values into commands.
package resp2parser

import (
	"fmt"
	"strconv"
	"strings"

	"kvaz/internal/command"
	"kvaz/internal/protocol/resp2"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

// Parse converts a RESP2 value into one or more commands.
//
// flat array -> one command.
// array of arrays -> a batch of commands.
func (p *Parser) Parse(value resp2.Value) ([]command.Command, error) {
	if value.Kind != resp2.KindArray {
		return nil, ParseErrorExpectedArray.format(value.Kind)
	}

	if len(value.Elements) == 0 {
		return nil, ParseErrorEmptyArray
	}

	// nested array -> a batch of commands.
	if value.Elements[0].Kind == resp2.KindArray {
		return p.parseBatch(value)
	}

	// the array represents one command
	cmd, err := p.parseCommand(value)
	if err != nil {
		return nil, err
	}

	return []command.Command{cmd}, nil
}

func (p *Parser) parseBatch(value resp2.Value) ([]command.Command, error) {
	commands := make([]command.Command, 0, len(value.Elements))

	for i, element := range value.Elements {
		if element.Kind != resp2.KindArray {
			return nil, ParseErrorInvalidBatchElement.format(i)
		}

		cmd, err := p.parseCommand(element)
		if err != nil {
			return nil, ParseErrorInvalidCommandAtBatch.format(i, err)
		}

		commands = append(commands, cmd)
	}

	return commands, nil
}

func (p *Parser) parseCommand(value resp2.Value) (command.Command, error) {
	if value.Kind != resp2.KindArray {
		return nil, ParseErrorExpectedCommandArray
	}

	if len(value.Elements) == 0 {
		return nil, ParseErrorEmptyArray
	}

	args := make([]string, len(value.Elements))

	for i, element := range value.Elements {
		arg, err := extractString(element)
		if err != nil {
			return nil, ParseErrorInvalidArgument.format(i, err)
		}

		args[i] = arg
	}

	name := strings.ToUpper(args[0])
	args = args[1:]

	switch name {
	case "PING":
		if err := validateArity(name, args, 0); err != nil {
			return nil, err
		}
		return command.NewPingCommand(), nil

	case "GET":
		if err := validateArity(name, args, 1); err != nil {
			return nil, err
		}
		return command.NewGetCommand(args[0]), nil

	case "SET":
		if err := validateArity(name, args, 2); err != nil {
			return nil, err
		}
		return command.NewSetCommand(args[0], args[1]), nil

	case "DEL":
		if err := validateArity(name, args, 1); err != nil {
			return nil, err
		}
		return command.NewDelCommand(args[0]), nil

	case "EXISTS":
		if err := validateArity(name, args, 1); err != nil {
			return nil, err
		}
		return command.NewExistsCommand(args[0]), nil

	case "EXPIRE":
		if err := validateArity(name, args, 2); err != nil {
			return nil, err
		}

		seconds, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, fmt.Errorf("invalid EXPIRE duration %q: %w", args[1], err)
		}

		return command.NewExpireCommand(args[0], int64(seconds)), nil

	case "TTL":
		if err := validateArity(name, args, 1); err != nil {
			return nil, err
		}
		return command.NewTTLCommand(args[0]), nil

	default:
		return nil, ParseErrorUnknownCommand
	}
}

func extractString(value resp2.Value) (string, error) {
	switch value.Kind {
	case resp2.KindSimpleString:
		return value.Text, nil

	case resp2.KindBulkString:
		return string(value.Bytes), nil

	default:
		return "", ParseErrorBadValue
	}
}

func validateArity(name string, args []string, expected int) error {
	if len(args) != expected {
		return ParseErrorArity.format(
			name,
			expected,
			len(args),
		)
	}

	return nil
}
