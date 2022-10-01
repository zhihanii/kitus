package kitus

import (
	"context"
	"github.com/zhihanii/gio"
)

type Client interface {
	Call(ctx context.Context, method string, req, resp interface{}) error
}

type client struct {
	ctx context.Context

	rt gio.RoundTripper
}

func NewClient(ctx context.Context) (Client, error) {
	var err error

	c := &client{
		ctx: ctx,
	}

	c.rt, err = gio.NewRoundTripper(ctx, &gio.RTOptions{})
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) Call(ctx context.Context, method string, req, resp interface{}) error {
	var (
		err error
	)

	err = c.rt.RoundTrip(
		func(conn gio.Conn) error {
			var (
				err error
			)

			buffer := gio.NewByteBuffer(conn)

			args, err := GetCodec().Marshal(req)
			if err != nil {
				return err
			}

			m := newMessage(method, args)
			data, _ := encode(m)

			if _, err = buffer.WriteBytes(data); err != nil {
				return err
			}

			return nil
		},
		func(conn gio.Conn) error {
			var err error

			buffer := gio.NewByteBuffer(conn)

			//读取消息长度
			n, err := buffer.ReadUint32()
			if err != nil {
				return err
			}

			data := make([]byte, n)
			//读取消息内容
			if _, err = buffer.ReadBytes(data); err != nil {
				return err
			}

			//返回的数据是怎样的
			if err = GetCodec().Unmarshal(data, resp); err != nil {
				return err
			}

			return nil
		},
	)

	return err
}
