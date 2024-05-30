package parser

import (
	"bufio"
	"redis/interface/reply"
	"redis/protocol"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOne(t *testing.T) {
	cases := []struct {
		input  string
		expect reply.Reply
	}{
		{
			"+OK\r\n",
			protocol.MakeSimpleStr("OK"),
		},
		{
			":1000\r\n",
			protocol.MakeInteger(1000),
		},
		{
			"$6\r\nfoobar\r\n",
			protocol.MakeBulkStr("foobar"),
		},
		{
			"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			protocol.MakeMultiRaw([]reply.Reply{
				protocol.MakeBulkStr("foo"),
				protocol.MakeBulkStr("bar"),
			}),
		},
	}

	for i, c := range cases {
		br := bufio.NewReader(strings.NewReader(c.input))
		actual := parseOne(br)
		assert.Equal(t, c.expect, actual, "case %d fail", i+1)
	}
}

func TestParse(t *testing.T) {
	input := "+OK\r\n:1000\r\n$1\r\nfoobar\r\n+OK\r\n"
	br := strings.NewReader(input)
	expect := []reply.Reply{
		protocol.MakeSimpleStr("OK"),
		protocol.MakeInteger(1000),
		protocol.MakeBulkStr("f"),
		protocol.MakeSimpleErr("invalid protocol"),
	}

	replyCh := Parse(br)
	actual := []reply.Reply{}
	for reply := range replyCh {
		actual = append(actual, reply)
	}

	assert.ElementsMatch(t, expect, actual)
}
