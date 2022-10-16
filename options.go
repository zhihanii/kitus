package kitus

type Option func(o *options)

type options struct {
	interceptor      ServerInterceptor
	numServerWorkers uint32
}

func Interceptor(i ServerInterceptor) Option {
	return func(o *options) {
		o.interceptor = i
	}
}

func NumServerWorkers(num uint32) Option {
	return func(o *options) {
		o.numServerWorkers = num
	}
}
