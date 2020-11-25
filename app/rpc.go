package app

import (
	"context"

	"github.com/gonitro/nitro/app/client"
	rpcClient "github.com/gonitro/nitro/app/client/rpc"
	mevent "github.com/gonitro/nitro/app/event/memory"
	sock "github.com/gonitro/nitro/app/network/socket"
	"github.com/gonitro/nitro/app/registry/memory"
	"github.com/gonitro/nitro/app/router/static"
	"github.com/gonitro/nitro/app/server"
	rpcServer "github.com/gonitro/nitro/app/server/rpc"
)

type rpcProgram struct {
	opts Options
}

func (s *rpcProgram) Name(name string) {
	s.opts.Server.Init(
		server.Name(name),
	)
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *rpcProgram) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}
}

func (s *rpcProgram) Options() Options {
	return s.opts
}

func (s *rpcProgram) Execute(name, ep string, req, rsp interface{}) error {
	r := s.Client().NewRequest(name, ep, req)
	return s.Client().Call(context.Background(), r, rsp)
}

func (s *rpcProgram) Broadcast(event string, msg interface{}) error {
	m := s.Client().NewMessage(event, msg)
	return s.Client().Publish(context.Background(), m)
}

func (s *rpcProgram) Register(v interface{}) error {
	h := s.Server().NewHandler(v)
	return s.Server().Handle(h)
}

func (s *rpcProgram) Subscribe(event string, v interface{}) error {
	sub := s.Server().NewSubscriber(event, v)
	return s.Server().Subscribe(sub)
}

func (s *rpcProgram) Client() client.Client {
	return s.opts.Client
}

func (s *rpcProgram) Server() server.Server {
	return s.opts.Server
}

func (s *rpcProgram) String() string {
	return "rpc"
}

func (s *rpcProgram) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *rpcProgram) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *rpcProgram) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	// wait on context cancel
	<-s.opts.Context.Done()

	return s.Stop()
}

// New returns a new application program
func New(opts ...Option) *rpcProgram {
	b := mevent.NewBroker()
	c := rpcClient.NewClient()
	s := rpcServer.NewServer()
	r := memory.NewRegistry()
	t := sock.NewTransport()
	st := static.NewRouter()

	// set client options
	c.Init(
		client.Router(st),
		client.Broker(b),
		client.Registry(r),
		client.Transport(t),
	)

	// set server options
	s.Init(
		server.Broker(b),
		server.Registry(r),
		server.Transport(t),
	)

	// define local opts
	options := Options{
		Broker:   b,
		Client:   c,
		Server:   s,
		Registry: r,
		Context:  context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &rpcProgram{
		opts: options,
	}
}
