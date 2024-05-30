package tcp

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTcpServerMacConn(t *testing.T) {
	addr := "127.0.0.1:"
	listener, err := net.Listen("tcp", addr)
	require.Nil(t, err)

	maxConn := 1
	handler := &EchoHandler{}
	go ListenAndServe(listener, maxConn, handler)

	conn, err := net.Dial("tcp", listener.Addr().String())
	require.Nil(t, err)

	conn2, err := net.Dial("tcp", listener.Addr().String())
	require.Nil(t, err)
	conn2.Write([]byte("hello\n"))
	conn2.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	_, err = conn2.Read(make([]byte, 5))
	// expect read timeout, cause maxCon is 1
	require.NotNil(t, err)

	conn.Close()

	conn2.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	_, err = conn2.Read(make([]byte, 5))
	require.Nil(t, err)
}
