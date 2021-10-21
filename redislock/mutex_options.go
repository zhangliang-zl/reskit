package redislock

import "time"

var (
	DefaultKeyPrefix     = "MutexLocker:"
	DefaultRetryInterval = 20 * time.Millisecond
	DefaultLocked        = 30 * time.Second
	DefaultLockWaiting   = 120 * time.Second
)

type Options struct {
	duration      time.Duration
	lockWaiting   time.Duration
	retryInterval time.Duration
	keyPrefix     string
}

type Option func(o *Options)

func WithKeyPrefix(prefix string) Option {
	return func(o *Options) {
		o.keyPrefix = prefix
	}
}

func WithRetryInterval(duration time.Duration) Option {
	return func(o *Options) {
		o.retryInterval = duration
	}
}

func WithLockWaiting(duration time.Duration) Option {
	return func(o *Options) {
		o.lockWaiting = duration
	}
}

func WithLockTime(duration time.Duration) Option {
	return func(o *Options) {
		o.duration = duration
	}
}
