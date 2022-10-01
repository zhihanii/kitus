package kitus

import (
	"context"
	"github.com/zhihanii/gio"
	"net"
)

type Transport interface {
	Serve(ln net.Listener) error
}

type transport struct {
	ctx context.Context

	eh gio.EventHandler
}

func newTransport(eh gio.EventHandler) Transport {
	return &transport{eh: eh}
}

func (t *transport) Serve(ln net.Listener) error {
	var err error
	tcpLn, err := gio.NewTCPListener(t.ctx, t.eh, &gio.Options{})
	if err != nil {
		return err
	}
	err = tcpLn.Serve(ln)
	if err != nil {
		return err
	}
	return nil
}
