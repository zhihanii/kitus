package kitus

type Options struct {
	interceptor      ServerInterceptor
	numServerWorkers uint32
}
