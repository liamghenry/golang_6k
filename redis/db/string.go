package db

import (
	"redis/interface/reply"
	"redis/protocol"
	"strconv"
	"strings"
	"time"
)

func init() {
	registerCMD("set", cmdSet)
	registerCMD("get", cmdGet)
	registerCMD("mset", cmdMSet)
}

// cmdSet handles the SET command.
func cmdSet(db *DB, args [][]byte) reply.Reply {
	if len(args) < 2 {
		return protocol.MakeSimpleErr("ERR wrong number of arguments for 'set' command")
	}

	ttlMillseconds := 0
	if len(args) > 2 {
		for i := 2; i < len(args); i++ {
			if strings.ToLower(string(args[i])) == "ex" {
				if i+1 < len(args) {
					seconds, err := strconv.Atoi(string(args[i+1]))
					if err != nil {
						return protocol.MakeSimpleErr("ERR value is not an integer or out of range")
					}
					ttlMillseconds = seconds * 1000
				}
			} else if strings.ToLower(string(args[i])) == "px" {
				if i+1 < len(args) {
					millseconds, err := strconv.Atoi(string(args[i+1]))
					if err != nil {
						return protocol.MakeSimpleErr("ERR value is not an integer or out of range")
					}
					ttlMillseconds = millseconds
				}
			}
		}
	}

	count := db.items.Set(string(args[0]), args[1])
	if ttlMillseconds > 0 {
		db.setTTL(string(args[0]), time.Now().Add(time.Duration(ttlMillseconds)*time.Millisecond))
	}
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

// cmdMset handles the MSET command.
func cmdMSet(db *DB, args [][]byte) reply.Reply {
	keys := []string{}
	values := []interface{}{}

	// args is key1, value1, key2, value2, ..., get keys values
	for i := 0; i < len(args); i += 2 {
		keys = append(keys, string(args[i]))
		values = append(values, args[i+1])
	}

	if len(keys) != len(values) {
		return protocol.MakeSimpleErr("ERR wrong number of arguments for 'mset' command")
	}

	db.items.SetMulti(keys, values)
	return protocol.MakeSimpleStr("OK")
}
