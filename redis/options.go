package redis

import "github.com/go-kratos/kratos/v2/log"

var DefaultOption = &Options{
	addr:     "127.0.0.1:6379",
	password: "",
	dbNum:    0,
	logger:   log.NewHelper(log.With(log.DefaultLogger, "tag", "redis")),
}

type Option func(options *Options)

func Logger(logger *log.Helper) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func Address(address string) Option {
	return func(o *Options) {
		o.addr = address
	}
}

func Password(pass string) Option {
	return func(o *Options) {
		o.password = pass
	}
}

func DbNum(dbNum int) Option {
	return func(o *Options) {
		o.dbNum = dbNum
	}
}
