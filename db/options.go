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
	DefaultDisableAutoPing = true
	DefaultLogger          = log.NewHelper(log.DefaultLogger, log.WithMessageKey("db"))
)

type Option func(options *Options)

func WithMaxOPenConns(num int) Option {
	return func(c *Options) {
		c.maxOpenConns = num
	}
}

func WithMaxIdleConns(num int) Option {
	return func(c *Options) {
		c.maxIdleConns = num
	}
}

func WithSlowThreshold(duration time.Duration) Option {
	return func(c *Options) {
		c.slowThreshold = duration
	}
}

func WithConnMaxLifetime(duration time.Duration) Option {
	return func(c *Options) {
		c.connMaxLifetime = duration
	}
}

func WithConnMaxIdleTime(duration time.Duration) Option {
	return func(c *Options) {
		c.connMaxIdleTime = duration
	}
}

func WithDisableAutoPing(close bool) Option {
	return func(c *Options) {
		c.disableAutoPing = close
	}
}

func WithLogger(logger *log.Helper) Option {
	return func(c *Options) {
		c.Logger = logger
	}
}
