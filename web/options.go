package web

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Options struct {
	address      string
	middlewares  []HandlerFunc
	readTimeout  time.Duration
	writeTimeout time.Duration
	logger       *log.Helper
}

var (
	DefaultAddress         = ":8080"
	DefaultReadTimeout     = time.Second * 300
	DefaultWriteTimeout    = time.Second * 300
	DefaultSlowThresholdMS = 200
	DefaultLogger          = log.NewHelper(log.With(log.DefaultLogger, "tag", "web"))
	DefaultMiddlewares     = []HandlerFunc{
		Recovery(DefaultLogger),
		Speed(DefaultLogger, DefaultSlowThresholdMS),
		Logging(DefaultLogger),
	}
)

type Option func(*Options)

func Address(address string) Option {
	return func(o *Options) {
		o.address = address
	}
}

func Middleware(middlewares ...HandlerFunc) Option {
	return func(o *Options) {
		o.middlewares = middlewares
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
