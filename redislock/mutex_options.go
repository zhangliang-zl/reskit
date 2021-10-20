package redislock

import "time"

const (
	DefaultKeyPrefix     = "MutexLocker:"
	DefaultRetryInterval = 50 * time.Millisecond
	DefaultLocked        = 10 * time.Second
	DefaultLockWaiting   = 30 * time.Second
)

type Options struct {
	Duration      time.Duration // default 10s
	LockWaiting   time.Duration // default 30s
	RetryInterval time.Duration // default 50ms
	KeyPrefix     string        // redis key prefix
}

type Option func(o *Options)

func KeyPrefix(prefix string) Option {
	return func(o *Options) {
		o.KeyPrefix = prefix
	}
}

func RetryInterval(duration time.Duration) Option {
	return func(o *Options) {
		o.RetryInterval = duration
	}
}
func LockWaiting(duration time.Duration) Option {
	return func(o *Options) {
		o.LockWaiting = duration
	}
}

func LockTime(duration time.Duration) Option {
	return func(o *Options) {
		o.Duration = duration
	}
}
