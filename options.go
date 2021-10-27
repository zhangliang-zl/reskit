package reskit

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs"
	"os"
)

type Option func(options *Options)

type Options struct {
	env  string
	name string

	ctx     context.Context
	servers []Server

	// hook functions
	beforeStart []func() error
	beforeStop  []func() error
	afterStart  []func() error
	afterStop   []func() error

	sigs   []os.Signal
	logger logs.Logger
}

func Name(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

func Env(env string) Option {
	return func(o *Options) {
		o.env = env
	}
}

func Servers(servers ...Server) Option {
	return func(o *Options) {
		o.servers = servers
	}
}

func Signal(sigs ...os.Signal) Option {
	return func(o *Options) {
		o.sigs = sigs
	}
}

func Logger(logger logs.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

func BeforeStart(f ...func() error) Option {
	return func(o *Options) {
		o.beforeStart = f
	}
}

func AfterStart(f ...func() error) Option {
	return func(o *Options) {
		o.afterStart = f
	}
}

func BeforeStop(f ...func() error) Option {
	return func(o *Options) {
		o.beforeStop = f
	}
}

func AfterStop(f ...func() error) Option {
	return func(o *Options) {
		o.afterStop = f
	}
}
