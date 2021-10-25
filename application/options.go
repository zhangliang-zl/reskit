package application

import (
	"github.com/go-kratos/kratos/v2/log"
	"os"
)

type Option func(options *Options)

type Options struct {
	env  string
	name string

	servers []Server

	// hook functions
	beforeStart []func() error
	beforeStop  []func() error
	afterStart  []func() error
	afterStop   []func() error

	sigs   []os.Signal
	logger *log.Helper
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

