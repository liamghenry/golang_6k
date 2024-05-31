package db

import (
	concurrentmap "redis/concurrent_map"
	"redis/interface/reply"
	"redis/protocol"
	"strings"
)

type DB struct {
	items *concurrentmap.ConcurrentMap
}

func NewDB() *DB {
	return &DB{
		items: concurrentmap.NewConcurrentMap(16),
	}
}

// Execute executes the command and returns the result.
func (db *DB) Execute(cmd string, args [][]byte) reply.Reply {
	if fn, ok := cmdTable[strings.ToLower(cmd)]; ok {
		return fn(db, args)
	}
	return protocol.MakeSimpleErr("unknown command " + cmd)
}