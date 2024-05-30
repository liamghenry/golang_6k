package main

import (
	"redis/tcp"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	err := tcp.ListenAndServe(":6379", &tcp.EchoHandler{})
	if err != nil {
		logrus.Error(err)
	}
}
