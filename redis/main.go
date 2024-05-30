package main

import (
	"redis/tcp"

	"github.com/sirupsen/logrus"
)

func main() {
	err := tcp.ListenAndServe(":6379", &tcp.EchoHandler{})
	logrus.Error(err)
}
