package mq

import "time"

type Options struct {
	minPriority    uint32
	maxPriority    uint32
	maxWorkingTTL  time.Duration
	failRetryDelay time.Duration
}

type Option func(o *Options)

func MinPriority(prior uint32) Option {
	return func(c *Options) {
		c.minPriority = prior
	}
}

func MaxPriority(prior uint32) Option {
	return func(c *Options) {
		c.maxPriority = prior
	}
}

func WorkingTTL(duration time.Duration) Option {
	return func(c *Options) {
		c.maxWorkingTTL = duration
	}
}

func FailRetryDelay(duration time.Duration) Option {
	return func(c *Options) {
		c.failRetryDelay = duration
	}
}
