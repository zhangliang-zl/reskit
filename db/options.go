package db

import (
	"github.com/zhangliang-zl/reskit/logs"
	"time"
)

type Option func(options *Options)

func MaxOpenConnections(num int) Option {
	return func(c *Options) {
		c.maxOpenConns = num
	}
}

func MaxIdleConnections(num int) Option {
	return func(c *Options) {
		c.maxIdleConns = num
	}
}

func SlowThreshold(duration time.Duration) Option {
	return func(c *Options) {
		c.slowThreshold = duration
	}
}

func ConnMaxLifetime(duration time.Duration) Option {
	return func(c *Options) {
		c.connMaxLifetime = duration
	}
}

func ConnMaxIdleTime(duration time.Duration) Option {
	return func(c *Options) {
		c.connMaxIdleTime = duration
	}
}

func DisableAutoPing(close bool) Option {
	return func(c *Options) {
		c.disableAutoPing = close
	}
}

func Logger(logger logs.Logger) Option {
	return func(c *Options) {
		c.Logger = logger
	}
}
