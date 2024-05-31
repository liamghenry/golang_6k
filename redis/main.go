package main

import (
	"net"
	"redis/redis"
	"redis/tcp"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	listener, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("server start at: ", listener.Addr())
	tcp.ListenAndServe(listener, 64, redis.MakeRedisServer())
}
