package kitus

import "context"

type ServiceInfo struct {
	ServiceName string

	HandlerType interface{}
	Handler     interface{}

	Methods map[string]*MethodInfo
}

type MethodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error)

type MethodInfo struct {
	MethodName string
	Handler    MethodHandler
}
