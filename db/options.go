package db

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

var (
	DefaultMaxOpenConns    = 100
	DefaultMaxIdleConns    = 100
	DefaultSlowThreshold   = 100 * time.Millisecond
	DefaultConnMaxLifetime = 300 * time.Second
	DefaultConnMaxIdleTime = 300 * time.Second
	DefaultLogger          = log.NewHelper(log.With(log.DefaultLogger, "tag", "db"))
)

type Option func(options *Options)

func MaxOPenConns(num int) Option {
	return func(c *Options) {
		c.maxOpenConns = num
	}
}

func MaxIdleConns(num int) Option {
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

func Logger(logger *log.Helper) Option {
	return func(c *Options) {
		c.Logger = logger
	}
}
