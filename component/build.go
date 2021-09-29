package component

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/component/cache"
	"github.com/zhangliang-zl/reskit/component/db"
	"github.com/zhangliang-zl/reskit/component/lock"
	kvstore "github.com/zhangliang-zl/reskit/component/redis"
	"github.com/zhangliang-zl/reskit/logs"
)

func NewRedis(opts kvstore.Options, logger logs.Logger) (Interface, error) {
	instance, err := kvstore.New(opts, logger)
	if err != nil {
		return nil, err
	}

	return Make(instance, instance.Close), nil
}

func NewMutexFactory(prefix string, kvstroe *redis.Client, logger logs.Logger) (Interface, error) {
	factory := lock.NewRedisMutexFactory(logger, kvstroe, prefix)
	instance := Make(factory, nil)
	return instance, nil
}

func NewDB(opts db.Options, logger logs.Logger) (Interface, error) {
	client, err := db.New(opts, logger)
	if err != nil {
		return nil, err
	}

	closeFunc := func() error {
		s, err := client.DB()
		if err != nil {
			return err
		}
		return s.Close()
	}

	return Make(client, closeFunc), nil
}

func NewRedisCache(prefix string, kvstroe *redis.Client, logger logs.Logger) Interface {
	instance := cache.NewRedisCache(kvstroe, logger, prefix)
	return Make(instance, nil)
}

func RegisterMemoryCache(size int) Interface {
	instance := cache.NewMemoryCache(size)
	return Make(instance, nil)
}
