package redislock

import (
	"github.com/zhangliang-zl/reskit/logs"
	"time"
)

type Options struct {
	duration      time.Duration
	lockWaiting   time.Duration
	retryInterval time.Duration
	keyPrefix     string
	logger        logs.Logger
}

type Option func(o *Options)

func KeyPrefix(prefix string) Option {
	return func(o *Options) {
		o.keyPrefix = prefix
	}
}

func RetryInterval(duration time.Duration) Option {
	return func(o *Options) {
		o.retryInterval = duration
	}
}

func LockWaiting(duration time.Duration) Option {
	return func(o *Options) {
		o.lockWaiting = duration
	}
}

func LockTime(duration time.Duration) Option {
	return func(o *Options) {
		o.duration = duration
	}
}
