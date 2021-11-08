package server

import (
	"github.com/zhangliang-zl/reskit/logs"
	"time"
)

type Options struct {
	address      string
	middleware   []HandlerFunc
	readTimeout  time.Duration
	writeTimeout time.Duration
	logger       logs.Logger
}

type Option func(*Options)

func Address(address string) Option {
	return func(o *Options) {
		o.address = address
	}
}

func Middleware(middleware ...HandlerFunc) Option {
	return func(o *Options) {
		o.middleware = middleware
	}
}

func ReadTimeout(duration time.Duration) Option {
	return func(o *Options) {
		o.readTimeout = duration
	}
}

func WriteTimeout(duration time.Duration) Option {
	return func(o *Options) {
		o.writeTimeout = duration
	}
}
