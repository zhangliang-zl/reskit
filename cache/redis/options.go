package redis

import "github.com/go-kratos/kratos/v2/log"

var (
	DefaultLogger    = log.NewHelper(log.With(log.DefaultLogger, "tag", "cache"))
	DefaultKeyPrefix = "Caching:"
)

type Option func(c *Cache)

func Prefix(key string) Option {
	return func(c *Cache) {
		c.keyPrefix = key
	}
}

func Logger(logger *log.Helper) Option {
	return func(c *Cache) {
		c.logger = logger
	}
}
