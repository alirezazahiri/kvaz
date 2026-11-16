package resp2parser_test

import (
	"reflect"
	"testing"

	"kvaz/internal/command"
	resp2parser "kvaz/internal/parser/resp2"
	"kvaz/internal/protocol/resp2"
)

func simple(s string) resp2.Value {
	return resp2.Value{
		Kind: resp2.KindSimpleString,
		Text: s,
	}
}

func bulk(s string) resp2.Value {
	return resp2.Value{
		Kind:  resp2.KindBulkString,
		Bytes: []byte(s),
	}
}

func array(values ...resp2.Value) resp2.Value {
	return resp2.Value{
		Kind:     resp2.KindArray,
		Elements: values,
		Integer:  int64(len(values)),
	}
}

func TestRESP2Parser_Parse_SingleCommands(t *testing.T) {
	p := resp2parser.New()

	tests := []struct {
		name  string
		input resp2.Value
		want  []command.Command
	}{
		{
			name:  "PING",
			input: array(simple("PING")),
			want:  []command.Command{command.NewPingCommand()},
		},
		{
			name:  "GET",
			input: array(simple("GET"), simple("foo")),
			want:  []command.Command{command.NewGetCommand("foo")},
		},
		{
			name: "SET",
			input: array(
				simple("SET"),
				simple("foo"),
				simple("bar"),
			),
			want: []command.Command{
				command.NewSetCommand("foo", "bar"),
			},
		},
		{
			name:  "DEL",
			input: array(simple("DEL"), simple("foo")),
			want:  []command.Command{command.NewDelCommand("foo")},
		},
		{
			name:  "EXISTS",
			input: array(simple("EXISTS"), simple("foo")),
			want:  []command.Command{command.NewExistsCommand("foo")},
		},
		{
			name: "EXPIRE",
			input: array(
				simple("EXPIRE"),
				simple("foo"),
				simple("60"),
			),
			want: []command.Command{
				command.NewExpireCommand("foo", 60),
			},
		},
		{
			name:  "TTL",
			input: array(simple("TTL"), simple("foo")),
			want:  []command.Command{command.NewTTLCommand("foo")},
		},
		{
			name:  "lowercase command",
			input: array(simple("get"), simple("foo")),
			want:  []command.Command{command.NewGetCommand("foo")},
		},
		{
			name:  "mixed case command",
			input: array(simple("SeT"), simple("foo"), simple("bar")),
			want: []command.Command{
				command.NewSetCommand("foo", "bar"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestRESP2Parser_Parse_BulkStringArguments(t *testing.T) {
	p := resp2parser.New()

	input := array(
		bulk("SET"),
		bulk("foo"),
		bulk("bar"),
	)

	want := []command.Command{
		command.NewSetCommand("foo", "bar"),
	}

	got, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}
}

func TestRESP2Parser_Parse_MixedStringTypes(t *testing.T) {
	p := resp2parser.New()

	input := array(
		simple("SET"),
		bulk("foo"),
		simple("bar"),
	)

	want := []command.Command{
		command.NewSetCommand("foo", "bar"),
	}

	got, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}
}

func TestRESP2Parser_Parse_MultipleCommands(t *testing.T) {
	p := resp2parser.New()

	input := array(
		array(simple("SET"), simple("foo"), simple("bar")),
		array(simple("GET"), simple("foo")),
		array(simple("PING")),
	)

	want := []command.Command{
		command.NewSetCommand("foo", "bar"),
		command.NewGetCommand("foo"),
		command.NewPingCommand(),
	}

	got, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}
}

func TestRESP2Parser_Parse_EmptyArray(t *testing.T) {
	p := resp2parser.New()

	got, err := p.Parse(array())
	if err == nil {
		t.Fatal("Parse() expected an error")
	}

	if got != nil {
		t.Errorf("Parse() commands = %#v, want nil", got)
	}
}

func TestRESP2Parser_Parse_InvalidTopLevelValues(t *testing.T) {
	p := resp2parser.New()

	tests := []struct {
		name  string
		input resp2.Value
	}{
		{"simple string", simple("PING")},
		{"bulk string", bulk("PING")},
		{"integer", resp2.Value{Kind: resp2.KindInteger, Integer: 42}},
		{"null", resp2.Null()},
		{"error", resp2.Value{Kind: resp2.KindError, Text: "ERR"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.input)

			if err == nil {
				t.Fatal("Parse() expected an error")
			}

			if got != nil {
				t.Errorf("Parse() commands = %#v, want nil", got)
			}
		})
	}
}

func TestRESP2Parser_Parse_InvalidCommandArguments(t *testing.T) {
	p := resp2parser.New()

	tests := []struct {
		name  string
		input resp2.Value
	}{
		{
			name: "integer command name",
			input: array(
				resp2.Value{Kind: resp2.KindInteger, Integer: 1},
			),
		},
		{
			name: "integer argument",
			input: array(
				simple("GET"),
				resp2.Value{Kind: resp2.KindInteger, Integer: 1},
			),
		},
		{
			name: "null argument",
			input: array(
				simple("GET"),
				resp2.Null(),
			),
		},
		{
			name: "nested argument",
			input: array(
				simple("GET"),
				array(simple("foo")),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.input)

			if err == nil {
				t.Fatal("Parse() expected an error")
			}

			if got != nil {
				t.Errorf("Parse() commands = %#v, want nil", got)
			}
		})
	}
}

func TestRESP2Parser_Parse_InvalidArity(t *testing.T) {
	p := resp2parser.New()

	tests := []struct {
		name  string
		input resp2.Value
	}{
		{"PING with argument", array(simple("PING"), simple("extra"))},
		{"GET without key", array(simple("GET"))},
		{"GET with extra argument", array(simple("GET"), simple("foo"), simple("bar"))},
		{"SET without value", array(simple("SET"), simple("foo"))},
		{"SET with extra argument", array(simple("SET"), simple("foo"), simple("bar"), simple("extra"))},
		{"DEL without key", array(simple("DEL"))},
		{"EXISTS without key", array(simple("EXISTS"))},
		{"EXPIRE without duration", array(simple("EXPIRE"), simple("foo"))},
		{"TTL without key", array(simple("TTL"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.input)

			if err == nil {
				t.Fatal("Parse() expected an error")
			}

			if got != nil {
				t.Errorf("Parse() commands = %#v, want nil", got)
			}
		})
	}
}

func TestRESP2Parser_Parse_UnknownCommand(t *testing.T) {
	p := resp2parser.New()

	got, err := p.Parse(array(simple("UNKNOWN")))

	if err == nil {
		t.Fatal("Parse() expected an error")
	}

	if got != nil {
		t.Errorf("Parse() commands = %#v, want nil", got)
	}
}

func TestRESP2Parser_Parse_InvalidExpireDuration(t *testing.T) {
	p := resp2parser.New()

	tests := []string{
		"abc",
		"1.5",
		"",
		"9223372036854775808",
	}

	for _, duration := range tests {
		t.Run(duration, func(t *testing.T) {
			input := array(
				simple("EXPIRE"),
				simple("foo"),
				simple(duration),
			)

			got, err := p.Parse(input)

			if err == nil {
				t.Fatal("Parse() expected an error")
			}

			if got != nil {
				t.Errorf("Parse() commands = %#v, want nil", got)
			}
		})
	}
}

func TestRESP2Parser_Parse_InvalidBatchIsAtomic(t *testing.T) {
	p := resp2parser.New()

	input := array(
		array(simple("PING")),
		array(simple("GET"), simple("foo")),
		array(simple("UNKNOWN")),
	)

	got, err := p.Parse(input)

	if err == nil {
		t.Fatal("Parse() expected an error")
	}

	if got != nil {
		t.Errorf("Parse() commands = %#v, want nil", got)
	}
}
