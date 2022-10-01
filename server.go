package kitus

import (
	"context"
	"encoding/binary"
	"github.com/zhihanii/gio"
	"net"
	"reflect"
	"strings"
	"sync/atomic"
)

type Server interface {
	RegisterService(svcInfo *ServiceInfo, handler interface{})
	Serve(ln net.Listener) error
}

type server struct {
	ctx context.Context

	opts *Options

	services map[string]*ServiceInfo

	msgChannels []chan *Message
}

func NewServer(ctx context.Context, opts *Options) Server {
	s := &server{
		ctx:      ctx,
		opts:     opts,
		services: make(map[string]*ServiceInfo),
	}
	s.initServerWorkers()
	return s
}

func (s *server) initServerWorkers() {
	s.msgChannels = make([]chan *Message, s.opts.numServerWorkers)
	for i := uint32(0); i < s.opts.numServerWorkers; i++ {
		s.msgChannels[i] = make(chan *Message)
		go s.serveChan(s.msgChannels[i])
	}
}

func (s *server) serveChan(ch chan *Message) {
	for {
		select {
		case m := <-ch:
			s.handleMessage(m)
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *server) handleMessage(m *Message) {
	sm := string(m.RpcInfo.Method)
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		//log
		return
	}

	service := sm[:pos]
	method := sm[pos+1:]

	svcInfo, knownService := s.services[service]
	if knownService {
		if md, ok := svcInfo.Methods[method]; ok {
			s.processRPC(svcInfo.Handler, md, m)
		}
	}
}

func (s *server) processRPC(srv interface{}, md *MethodInfo, m *Message) {
	dec := func(v interface{}) error {
		if err := GetCodec().Unmarshal(m.Args, v); err != nil {
			return err
		}
		return nil
	}
	resp, err := md.Handler(srv, s.ctx, dec)
	if err != nil {

	}
	err = s.sendResponse(m.Writer, resp)
	if err != nil {

	}
}

func (s *server) sendResponse(w gio.Writer, msg interface{}) error {
	data, err := GetCodec().Marshal(msg)
	if err != nil {
		return err
	}
	b := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(b[:4], uint32(len(data)))
	_ = copy(b[4:], data)
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) RegisterService(svcInfo *ServiceInfo, handler interface{}) {
	if handler != nil {
		ht := reflect.TypeOf(svcInfo.HandlerType).Elem()
		h := reflect.TypeOf(handler)
		if !h.Implements(ht) {
			//fatal
		}
	}
	svcInfo.Handler = handler
	s.services[svcInfo.ServiceName] = svcInfo
}

func (s *server) Serve(ln net.Listener) error {
	var err error
	var roundRobinCounter uint32
	handle := func(m *Message) {
		select {
		case s.msgChannels[atomic.AddUint32(&roundRobinCounter, 1)%s.opts.numServerWorkers] <- m:
		}
	}
	sh := &srvHandler{Handle: handle}
	t := newTransport(sh)
	err = t.Serve(ln)
	if err != nil {
		return err
	}
	return nil
}
