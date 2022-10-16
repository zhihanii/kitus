package kitus

import "context"

type ServerInfo struct {
	Server     interface{}
	FullMethod string
}

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

type ServerInterceptor func(ctx context.Context, req interface{}, info *ServerInfo, handler Handler) (resp interface{}, err error)
