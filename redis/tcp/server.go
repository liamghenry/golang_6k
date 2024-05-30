package tcp

import (
	"context"
	"net"
	"redis/interface/tcp"

	"github.com/sirupsen/logrus"
)

func ListenAndServe(addr string, handler tcp.Handler) error {
	// start a tcp server
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logrus.Info("tcp server started on ", addr)

	// accept tcp connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Error("accept error: ", err)
			continue
		}

		// handle tcp connection
		ctx := context.Background()
		go handler.Handle(ctx, conn)
	}
}
