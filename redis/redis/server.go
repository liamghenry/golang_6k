package redis

import (
	"context"
	"net"
	"redis/parser"
	"redis/protocol"
	"sync"

	"github.com/sirupsen/logrus"
)

type Server struct {
	activeConn sync.Map
}

func MakeRedisServer() *Server {
	return &Server{}
}

// TODO 考虑 ctx timeout 情况
func (s *Server) Handle(ctx context.Context, conn net.Conn) {
	s.activeConn.Store(conn, struct{}{})

	requestCh := parser.Parse(conn)
	for payload := range requestCh {
		if _, ok := payload.(*protocol.SimpleErr); ok {
			conn.Write(protocol.MakeSimpleErr("invalid protocol").Marshal())
			conn.Close()
			s.activeConn.Delete(conn)
			return
		}
		logrus.Debugf("receive payload: %v\n", payload)
		conn.Write(payload.Marshal())
	}
}

func (s *Server) Close() error {
	s.activeConn.Range(func(key, value interface{}) bool {
		conn := key.(net.Conn)
		conn.Close()
		return true
	})
	return nil
}
