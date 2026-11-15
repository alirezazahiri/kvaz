package resp2_test

import (
	"reflect"
	"testing"

	"kvaz/internal/protocol/resp2"
)

func TestEncode_MatchesDecoderRoundTrip(t *testing.T) {
	inputs := []struct {
		wire string
		val  resp2.Value
	}{
		{"+OK\r\n", resp2.SimpleString("OK")},
		{"-ERR unknown command\r\n", resp2.Error("ERR unknown command")},
		{":1000\r\n", resp2.Integer(1000)},
		{"$-1\r\n", resp2.Null()},
		{
			"$6\r\nfoobar\r\n",
			resp2.Value{Kind: resp2.KindBulkString, Bytes: []byte("foobar")},
		},
		{
			"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			resp2.Array([]resp2.Value{
				{Kind: resp2.KindBulkString, Bytes: []byte("foo")},
				{Kind: resp2.KindBulkString, Bytes: []byte("bar")},
			}),
		},
	}

	for _, tt := range inputs {
		t.Run(tt.wire, func(t *testing.T) {
			got, err := resp2.Encode(tt.val)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			if string(got) != tt.wire {
				t.Fatalf("got %q, want %q", got, tt.wire)
			}

			p := resp2.New()
			p.Feed(got)
			decoded, ok, err := p.Next()
			if err != nil || !ok {
				t.Fatalf("decode: ok=%v err=%v", ok, err)
			}
			if !reflect.DeepEqual(decoded, tt.val) {
				t.Fatalf("round-trip: got %#v, want %#v", decoded, tt.val)
			}
		})
	}
}