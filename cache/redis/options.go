package redis

import "github.com/go-kratos/kratos/v2/log"

var (
	DefaultLogger    = log.NewHelper(log.DefaultLogger, log.WithMessageKey("redis"))
	DefaultKeyPrefix = "Caching:"
)

type Option func(c *Cache)

func WithPrefix(key string) Option {
	return func(c *Cache) {
		c.keyPrefix = key
	}
}

func WithLogger(logger *log.Helper) Option {
	return func(c *Cache) {
		c.logger = logger
	}
}
