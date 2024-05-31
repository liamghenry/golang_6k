package db

import (
	"redis/interface/reply"
	"redis/protocol"
)

func init() {
	registerCMD("set", cmdSet)
	registerCMD("get", cmdGet)
}

// cmdSet handles the SET command.
func cmdSet(db *DB, args [][]byte) reply.Reply {
	if len(args) < 2 {
		return protocol.MakeSimpleErr("ERR wrong number of arguments for 'set' command")
	}
	count := db.items.Set(string(args[0]), args[1])
	if count > 0 {
		return protocol.MakeSimpleStr("OK")
	}
	return protocol.MakeNil()
}

// cmdGet handles the GET command.
func cmdGet(db *DB, args [][]byte) reply.Reply {
	if len(args) != 1 {
		return protocol.MakeSimpleErr("ERR wrong number of arguments for 'get' command")
	}
	v, ok := db.items.Get(string(args[0]))
	if !ok {
		return protocol.MakeNil()
	}
	// NOTE
	a := string(v.([]byte))
	return protocol.MakeBulkStr(a)
}
