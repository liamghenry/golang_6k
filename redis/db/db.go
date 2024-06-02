package db

import (
	concurrentmap "redis/concurrent_map"
	"redis/interface/reply"
	"redis/protocol"
	"redis/timewheel"
	"strings"
	"time"
)

type DB struct {
	items  *concurrentmap.ConcurrentMap
	ttlMap *concurrentmap.ConcurrentMap
}

func NewDB() *DB {
	return &DB{
		// TODO 调整大小
		items:  concurrentmap.NewConcurrentMap(16),
		ttlMap: concurrentmap.NewConcurrentMap(16),
	}
}

// Execute executes the command and returns the result.
func (db *DB) Execute(cmd string, args [][]byte) reply.Reply {
	if fn, ok := cmdTable[strings.ToLower(cmd)]; ok {
		return fn(db, args)
	}
	return protocol.MakeSimpleErr("unknown command " + cmd)
}

// setTTL sets the ttl for the key.
func (db *DB) setTTL(key string, ttl time.Time) {
	db.ttlMap.Set(key, ttl)
	timewheel.At(key, ttl, func() {
		db.items.Remove(key)
		db.ttlMap.Remove(key)
	})
}
