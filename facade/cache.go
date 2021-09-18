package facade

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/cache"
	"github.com/zhangliang-zl/reskit/component"
)

func RegisterRedisCache(id string, kvstroe *redis.Client) error {
	logger := Logger(compCache)
	instance := component.Make(cache.NewRedisCache(kvstroe, logger, ""), nil, nil)
	return reskit.App().SetComponent(compCache, id, instance)
}

func RegisterMemoryCache(id string, size int) error {
	instance := component.Make(cache.NewMemoryCache(size), nil, nil)
	return reskit.App().SetComponent(compCache, id, instance)
}

func Cache(id string) cache.Cache {
	instance, ok := reskit.App().Component(compCache, id)
	if !ok {
		panic(compCache + noRegister)
	}

	return instance.(cache.Cache)
}
