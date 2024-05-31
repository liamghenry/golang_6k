package db

import "redis/interface/reply"

type ExecFunc func(db *DB, args [][]byte) reply.Reply

var cmdTable = map[string]ExecFunc{}

func registerCMD(cmd string, fn ExecFunc) {
	cmdTable[cmd] = fn
}