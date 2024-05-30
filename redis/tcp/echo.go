package tcp

import (
	"bufio"
	"context"
	"net"

	"github.com/sirupsen/logrus"
)

// implement tcp.Handler, response same message to client
type EchoHandler struct{}

// Handle implement tcp.Handler
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	for {
		// read one line, then response the same message
		bytes, err := br.ReadBytes('\n')
		if err != nil {
			logrus.Error("read error: ", err)
			return
		}
		conn.Write(bytes)
	}
}

// Close implement tcp.Handler
func (h *EchoHandler) Close() error {
	logrus.Info("close echo handler")
	return nil
}
