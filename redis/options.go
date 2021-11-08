package redis

import (
	"github.com/zhangliang-zl/reskit/logs"
)

type Option func(options *Options)

func Logger(logger logs.Logger) Option {
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
