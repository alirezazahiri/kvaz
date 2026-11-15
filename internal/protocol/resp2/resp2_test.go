package resp2_test

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"

	"kvaz/internal/protocol/resp2"
)

func TestNext_DecodesRESPValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  resp2.Value
	}{
		{"simple string", "+OK\r\n", resp2.SimpleString("OK")},
		{"simple string with spaces", "+hello world\r\n", resp2.SimpleString("hello world")},
		{"error", "-ERR unknown command\r\n", resp2.Error("ERR unknown command")},
		{"integer", ":1000\r\n", resp2.Integer(1000)},
		{"negative integer", ":-42\r\n", resp2.Integer(-42)},
		{"zero integer", ":0\r\n", resp2.Integer(0)},
		{
			"bulk string",
			"$6\r\nfoobar\r\n",
			resp2.Value{Kind: resp2.KindBulkString, Bytes: []byte("foobar")},
		},
		{
			"empty bulk string",
			"$0\r\n\r\n",
			resp2.Value{Kind: resp2.KindBulkString, Bytes: []byte{}},
		},
		{
			"bulk string containing CRLF",
			"$5\r\na\r\nbc\r\n",
			resp2.Value{Kind: resp2.KindBulkString, Bytes: []byte("a\r\nbc")},
		},
		{"null bulk string", "$-1\r\n", resp2.Value{Kind: resp2.KindNull}},
		{"empty array", "*0\r\n", resp2.Array([]resp2.Value{})},
		{"null array", "*-1\r\n", resp2.Value{Kind: resp2.KindNull}},
		{
			"array of bulk strings",
			"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			resp2.Array([]resp2.Value{
				{Kind: resp2.KindBulkString, Bytes: []byte("foo")},
				{Kind: resp2.KindBulkString, Bytes: []byte("bar")},
			}),
		},
		{
			"nested array",
			"*2\r\n:1\r\n*2\r\n+hello\r\n:42\r\n",
			resp2.Array([]resp2.Value{
				resp2.Integer(1),
				resp2.Array([]resp2.Value{
					resp2.SimpleString("hello"),
					resp2.Integer(42),
				}),
			}),
		},
		{
			"array containing null",
			"*2\r\n$-1\r\n$3\r\nfoo\r\n",
			resp2.Array([]resp2.Value{
				{Kind: resp2.KindNull},
				{Kind: resp2.KindBulkString, Bytes: []byte("foo")},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := resp2.New()
			p.Feed([]byte(tt.input))

			got, success, err := p.Next()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !success {
				t.Fatal("expected successful decoding")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNext_ReturnsInsufficientInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"simple string without terminator", "+OK"},
		{"simple string with partial terminator", "+OK\r"},
		{"integer without terminator", ":42"},
		{"bulk string without payload", "$5\r\n"},
		{"bulk string with partial payload", "$5\r\nhel"},
		{"bulk string with partial terminator", "$3\r\nfoo\r"},
		{"array without elements", "*2\r\n"},
		{"array with one missing element", "*2\r\n$3\r\nfoo\r\n"},
		{"nested array incomplete", "*1\r\n*1\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := resp2.New()
			p.Feed([]byte(tt.input))

			_, success, err := p.Next()

			if success {
				t.Fatal("expected unsuccessful decoding")
			}
			if err != resp2.ProtocolErrorInsufficientInputBytes {
				t.Fatalf("got error %v, want insufficient input", err)
			}
		})
	}
}

func TestNext_ReturnsMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unknown prefix", "?hello\r\n"},
		{"invalid integer", ":abc\r\n"},
		{"integer with trailing characters", ":12abc\r\n"},
		{"invalid bulk string length", "$abc\r\n"},
		{"invalid negative bulk length", "$-2\r\n"},
		{"invalid array length", "*abc\r\n"},
		{"invalid negative array length", "*-2\r\n"},
		{"bulk string missing terminator", "$1\r\nxXX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := resp2.New()
			p.Feed([]byte(tt.input))

			_, success, err := p.Next()

			if success {
				t.Fatal("expected unsuccessful decoding")
			}
			if err != resp2.ProtocolErrorInputMalformed {
				t.Fatalf("got error %v, want malformed input", err)
			}
		})
	}
}

func TestNext_DecodesFragmentedInput(t *testing.T) {
	input := []byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")

	want := resp2.Array([]resp2.Value{
		{Kind: resp2.KindBulkString, Bytes: []byte("GET")},
		{Kind: resp2.KindBulkString, Bytes: []byte("foo")},
	})

	for split := 1; split < len(input); split++ {
		t.Run(strconv.Itoa(split), func(t *testing.T) {
			p := resp2.New()

			p.Feed(input[:split])

			_, success, err := p.Next()
			if success {
				t.Fatal("message should not be complete yet")
			}
			if err != resp2.ProtocolErrorInsufficientInputBytes {
				t.Fatalf("expected insufficient input, got %v", err)
			}

			p.Feed(input[split:])

			got, success, err := p.Next()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !success {
				t.Fatal("expected successful decoding")
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %#v, want %#v", got, want)
			}
		})
	}
}

func TestNext_DecodesOneByteAtATime(t *testing.T) {
	input := []byte("$5\r\nhello\r\n")
	p := resp2.New()

	for i, b := range input {
		p.Feed([]byte{b})

		got, success, err := p.Next()

		if i < len(input)-1 {
			if success {
				t.Fatal("decoded before receiving complete input")
			}
			if err != resp2.ProtocolErrorInsufficientInputBytes {
				t.Fatalf("byte %d: expected insufficient input, got %v", i, err)
			}
			continue
		}

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !success {
			t.Fatal("expected successful decoding")
		}

		want := resp2.Value{
			Kind:  resp2.KindBulkString,
			Bytes: []byte("hello"),
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	}
}

func TestNext_DecodesMultipleValues(t *testing.T) {
	p := resp2.New()

	p.Feed([]byte("+PONG\r\n:42\r\n$3\r\nfoo\r\n"))

	expected := []resp2.Value{
		resp2.SimpleString("PONG"),
		resp2.Integer(42),
		{Kind: resp2.KindBulkString, Bytes: []byte("foo")},
	}

	for i, want := range expected {
		got, success, err := p.Next()

		if err != nil || !success {
			t.Fatalf("value %d: success=%v err=%v", i, success, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("value %d: got %#v, want %#v", i, got, want)
		}
	}

	_, success, err := p.Next()
	if success {
		t.Fatal("expected no additional complete values")
	}
	if err != resp2.ProtocolErrorInsufficientInputBytes {
		t.Fatalf("expected insufficient input, got %v", err)
	}
}

func TestFeed_DoesNotModifyPreviouslyFedData(t *testing.T) {
	p := resp2.New()

	input := []byte("+OK\r\n")
	p.Feed(input)

	copy(input, []byte("xxxxx"))

	got, success, err := p.Next()

	if err != nil || !success {
		t.Fatalf("success=%v err=%v", success, err)
	}
	if !reflect.DeepEqual(got, resp2.SimpleString("OK")) {
		t.Fatalf("unexpected value: %#v", got)
	}
}

func TestNext_PreservesRemainingBytes(t *testing.T) {
	p := resp2.New()

	p.Feed(bytes.Join([][]byte{
		[]byte("+OK\r\n"),
		[]byte(":123\r\n"),
	}, nil))

	expected := []resp2.Value{
		resp2.SimpleString("OK"),
		resp2.Integer(123),
	}

	for i, want := range expected {
		got, success, err := p.Next()

		if err != nil || !success {
			t.Fatalf("value %d: success=%v err=%v", i, success, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("value %d: got %#v, want %#v", i, got, want)
		}
	}
}
