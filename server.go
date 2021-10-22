package reskit

import "context"

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Get() interface{}
}

type server struct {
	ref   interface{}
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
	return s.ref
}

func NewService(ref interface{}, start, stop func(ctx context.Context) error) Server {
	return &server{
		ref:   ref,
		start: start,
		stop:  stop,
	}
}
