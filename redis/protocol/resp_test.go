package protocol

import (
	"redis/interface/reply"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocol(t *testing.T) {
	// test MultiRaw
	mr := MakeMultiRaw([]reply.Reply{MakeSimpleStr("OK"), MakeInteger(1)})
	expect := []byte("*2\r\n+OK\r\n:1\r\n")

	actual := mr.Marshal()

	assert.Equal(t, expect, actual)
}
