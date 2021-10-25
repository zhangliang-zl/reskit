package reskit

import "context"

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Get() interface{}
}

type server struct {
	obj   interface{}
	start func(ctx context.Context) error
	stop  func(ctx context.Context) error
}

func (s *server) Start(ctx context.Context) error {
	return s.start(ctx)
}

func (s *server) Stop(ctx context.Context) error {
	return s.stop(ctx)
}

func (s *server) Get() interface{} {
	return s.obj
}

func BuildServer(obj interface{}, start, stop func(ctx context.Context) error) Server {
	return &server{
		obj:   obj,
		start: start,
		stop:  stop,
	}
}
