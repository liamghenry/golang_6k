package protocol

import (
	"redis/interface/reply"
	"strconv"
)

const CRLF = "\r\n"

type SimpleStr struct {
	str string
}

func MakeSimpleStr(str string) *SimpleStr {
	return &SimpleStr{str: str}
}

func (s *SimpleStr) Marshal() []byte {
	return []byte("+" + s.str + CRLF)
}

// SimpleErr
type SimpleErr struct {
	err string
}

func MakeSimpleErr(err string) *SimpleErr {
	return &SimpleErr{err: err}
}

func (s *SimpleErr) Marshal() []byte {
	return []byte("-" + s.err + CRLF)
}

// Integer
type Integer struct {
	num int
}

func MakeInteger(num int) *Integer {
	return &Integer{num: num}
}

func (i *Integer) Marshal() []byte {
	return []byte(":" + strconv.Itoa(i.num) + CRLF)
}

// BulkStr
type BulkStr struct {
	str string
}

func MakeBulkStr(str string) *BulkStr {
	return &BulkStr{str: str}
}

func (b *BulkStr) Marshal() []byte {
	return []byte("$" + strconv.Itoa(len(b.str)) + CRLF + b.str + CRLF)
}

// MultiBulk
type MultiBulk struct {
	items [][]byte
}

func MakeMultiBulk(items [][]byte) *MultiBulk {
	return &MultiBulk{items: items}
}

func (m *MultiBulk) Marshal() []byte {
	buf := []byte("*" + strconv.Itoa(len(m.items)) + CRLF)
	for _, item := range m.items {
		buf = append(buf, MakeBulkStr(string(item)).Marshal()...)
	}
	return buf
}

// MultiRaw
type MultiRaw struct {
	items []reply.Reply
}

func MakeMultiRaw(items []reply.Reply) *MultiRaw {
	return &MultiRaw{items: items}
}

func (m *MultiRaw) Marshal() []byte {
	buf := []byte("*" + strconv.Itoa(len(m.items)) + CRLF)
	for _, item := range m.items {
		buf = append(buf, item.Marshal()...)
	}
	return buf
}

// Null
type Null struct{}

func MakeNull() *Null {
	return &Null{}
}

func (n *Null) Marshal() []byte {
	return []byte("$-1" + CRLF)
}
