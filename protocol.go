package kitus

import (
	"encoding/binary"
	"github.com/zhihanii/gio"
	"sync/atomic"
)

var (
	globalSeqID uint32
)

func getSeqID() uint32 {
	return atomic.AddUint32(&globalSeqID, 1)
}

type RpcInfo struct {
	SeqID  uint32
	Method []byte
}

type Message struct {
	RpcInfo RpcInfo
	Args    []byte
	Writer  gio.Writer
}

func newMessage(method string, args []byte) *Message {
	return &Message{
		RpcInfo: RpcInfo{
			SeqID:  getSeqID(),
			Method: []byte(method),
		},
		Args: args,
	}
}

func encode(m *Message) ([]byte, error) {
	mLen, aLen := len(m.RpcInfo.Method), len(m.Args)
	b := make([]byte, 10+mLen+aLen)
	binary.BigEndian.PutUint32(b[:4], m.RpcInfo.SeqID)
	binary.BigEndian.PutUint16(b[4:6], uint16(mLen))
	_ = copy(b[6:6+mLen], m.RpcInfo.Method)
	binary.BigEndian.PutUint32(b[6+mLen:10+mLen], uint32(aLen))
	_ = copy(b[10+mLen:], m.Args)
	return b, nil
}

func parse(data []byte) (*Message, error) {
	var (
		rpcInfo RpcInfo
	)
	m := new(Message)
	rpcInfo.SeqID = binary.BigEndian.Uint32(data[:4])
	mLen := binary.BigEndian.Uint16(data[4:6])
	rpcInfo.Method = make([]byte, mLen)
	_ = copy(rpcInfo.Method, data[6:6+mLen])
	m.RpcInfo = rpcInfo
	aLen := binary.BigEndian.Uint32(data[6+mLen : 10+mLen])
	m.Args = make([]byte, aLen)
	_ = copy(m.Args, data[10+mLen:])
	return m, nil
}
