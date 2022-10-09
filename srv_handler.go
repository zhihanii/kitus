package kitus

import (
	"context"
	"github.com/zhihanii/gio"
)

type srvHandler struct {
	Handle func(m *Message)
}

func (h *srvHandler) OnConnect(ctx context.Context, conn gio.Conn) error {
	return nil
}

func (h *srvHandler) OnRead(ctx context.Context, conn gio.Conn) error {
	var (
		length int32
		n      int
		err    error
	)
	length, err = readHeader(conn)
	if err != nil {
		return err
	}

	data := make([]byte, length)
	n, err = conn.Read(data)
	if err != nil || int32(n) != length {

	}

	m, err := parse(data)
	if err != nil {
		return err
	}

	m.Writer = conn

	h.Handle(m)

	return nil
}
