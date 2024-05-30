package main

import (
	"net"
	"redis/tcp"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	listener, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		logrus.Error(err)
	}
	tcp.ListenAndServe(listener, 64, &tcp.EchoHandler{})
}
