package service

import "context"

type Interface interface {
	Object() interface{}
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type buildIns struct {
	object    interface{}
	startFunc func(ctx context.Context) error
	stopFunc  func(ctx context.Context) error
}

func (c buildIns) Object() interface{} {
	return c.object
}

func (c buildIns) Stop(ctx context.Context) error {
	return c.stopFunc(ctx)
}

func (c buildIns) Start(ctx context.Context) error {
	return c.startFunc(ctx)
}

func Make(object interface{}, start, stop func(ctx context.Context) error) Interface {
	return buildIns{
		object:    object,
		startFunc: start,
		stopFunc:  stop,
	}
}
