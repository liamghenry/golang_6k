package db

import (
	"redis/protocol"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdSetAndGet(t *testing.T) {
	db := NewDB()

	reply := db.Execute("set", [][]byte{[]byte("foo"), []byte("bar")})
	assert.Equal(t, protocol.MakeNil().Marshal(), reply.Marshal())

	reply = db.Execute("get", [][]byte{[]byte("foo")})
	assert.Equal(t, protocol.MakeBulkStr("bar").Marshal(), reply.Marshal())

	reply = db.Execute("set", [][]byte{[]byte("foo"), []byte("bar02")})
	assert.Equal(t, protocol.MakeSimpleStr("OK").Marshal(), reply.Marshal())

	reply = db.Execute("get", [][]byte{[]byte("foo")})
	assert.Equal(t, protocol.MakeBulkStr("bar02").Marshal(), reply.Marshal())
}
