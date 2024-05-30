package parser

import (
	"bufio"
	"io"
	"redis/interface/reply"
	"redis/protocol"
	"strconv"
)

func Parse(reader io.Reader) <-chan reply.Reply {
	replyCh := make(chan reply.Reply)
	br := bufio.NewReader(reader)

	go func() {
		for {
			reply := parseOne(br)
			replyCh <- reply
			if _, ok := reply.(*protocol.SimpleErr); ok {
				close(replyCh)
				return
			}
		}

	}()

	return replyCh
}

func parseOne(br *bufio.Reader) reply.Reply {
	firstByte := make([]byte, 1)
	_, err := br.Read(firstByte)
	if err != nil {
		return protocol.MakeSimpleErr(err.Error())
	}
	switch firstByte[0] {
	// Simple String
	case '+':
		line, err := br.ReadString('\n')
		if err != nil {
			return protocol.MakeSimpleErr(err.Error())
		} else {
			return protocol.MakeSimpleStr(line[:len(line)-2])
		}
	// Integer
	case ':':
		line, err := br.ReadString('\n')
		if err != nil {
			return protocol.MakeSimpleErr(err.Error())
		} else {
			num, _ := strconv.Atoi(line[:len(line)-2])
			return protocol.MakeInteger(num)
		}
	// BulkStr
	case '$':
		line, err := br.ReadString('\n')
		if err != nil {
			return protocol.MakeSimpleErr(err.Error())
		} else {
			n, _ := strconv.Atoi(line[:len(line)-2])
			if n == -1 {
				return protocol.MakeBulkStr("")
			} else {
				buf := make([]byte, n+2)
				_, err := io.ReadFull(br, buf)
				if err != nil {
					return protocol.MakeSimpleErr(err.Error())
				} else {
					return protocol.MakeBulkStr(string(buf[:n]))
				}
			}
		}
	// MultiRaw
	case '*':
		line, err := br.ReadString('\n')
		if err != nil {
			return protocol.MakeSimpleErr(err.Error())
		} else {
			n, _ := strconv.Atoi(line[:len(line)-2])
			if n == -1 {
				return protocol.MakeMultiBulk(nil)
			} else {
				replies := make([][]byte, n)
				for i := 0; i < n; i++ {
					bulk, ok := parseOne(br).(*protocol.BulkStr)
					if !ok {
						return protocol.MakeSimpleErr("invalid protocol")
					}
					replies[i] = bulk.Marshal()
				}
				return protocol.MakeMultiBulk(replies)
			}
		}
	default:
		return protocol.MakeSimpleErr("invalid protocol")
	}
}
