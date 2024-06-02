package db

import (
	"redis/protocol"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
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

func TestCmdSetWithExpire(t *testing.T) {
	db := NewDB()

	reply := db.Execute("set", [][]byte{[]byte("foo"), []byte("bar"), []byte("EX"), []byte("1")})
	assert.Equal(t, protocol.MakeSimpleStr("OK").Marshal(), reply.Marshal())

	time.Sleep(time.Millisecond * 500)
	reply = db.Execute("get", [][]byte{[]byte("foo")})
	assert.Equal(t, protocol.MakeBulkStr("bar").Marshal(), reply.Marshal())

	time.Sleep(time.Millisecond * 500)
	reply = db.Execute("get", [][]byte{[]byte("foo")})
	assert.Equal(t, protocol.MakeNil().Marshal(), reply.Marshal())
}

func TestCmdSetWithPxExpire(t *testing.T) {
	db := NewDB()

	reply := db.Execute("set", [][]byte{[]byte("foo"), []byte("bar"), []byte("PX"), []byte("1200")})
	assert.Equal(t, protocol.MakeSimpleStr("OK").Marshal(), reply.Marshal())

	time.Sleep(time.Millisecond * 500)
	reply = db.Execute("get", [][]byte{[]byte("foo")})
	assert.Equal(t, protocol.MakeBulkStr("bar").Marshal(), reply.Marshal())

	// 因为 tw 精度是 1s，所以 1200ms 需要等 2s 才能被销毁
	time.Sleep(time.Millisecond * 1500)
	reply = db.Execute("get", [][]byte{[]byte("foo")})
	logrus.Info(string(reply.Marshal()))
	assert.Equal(t, protocol.MakeNil().Marshal(), reply.Marshal())
}
