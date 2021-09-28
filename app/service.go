package app

import "context"

type Service interface {
	Object() interface{}
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type serviceIns struct {
	object interface{}
	start  func(ctx context.Context) error
	stop   func(ctx context.Context) error
}

func (c serviceIns) Object() interface{} {
	return c.object
}

func (c serviceIns) Stop(ctx context.Context) error {
	return c.Stop(ctx)
}

func (c serviceIns) Start(ctx context.Context) error {
	return c.Start(ctx)
}

func MakeService( object interface{}, start, stop func(ctx context.Context) error) Service {
	return serviceIns{
		object: object,
		start:  start,
		stop:   stop,
	}
}
