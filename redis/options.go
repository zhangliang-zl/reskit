package redis

import "github.com/go-kratos/kratos/v2/log"

var DefaultOption = &Options{
	addr:     "127.0.0.1:6379",
	password: "",
	dbNum:    0,
	logger:   log.NewHelper(log.DefaultLogger, log.WithMessageKey("redis")),
}

type Option func(options *Options)

func WithLogger(logger *log.Helper) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func WithAddress(address string) Option {
	return func(o *Options) {
		o.addr = address
	}
}

func WithPassword(pass string) Option {
	return func(o *Options) {
		o.password = pass
	}
}

func WithDbNum(dbNum int) Option {
	return func(o *Options) {
		o.dbNum = dbNum
	}
}
