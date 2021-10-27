package redis

import (
	"github.com/zhangliang-zl/reskit/logs"
)

var (
	DefaultLogger    = logs.DefaultLogger("cache")
	DefaultKeyPrefix = "Caching:"
)

type Option func(c *Cache)

func Prefix(key string) Option {
	return func(c *Cache) {
		c.keyPrefix = key
	}
}

func Logger(logger logs.Logger) Option {
	return func(c *Cache) {
		c.logger = logger
	}
}
