package facade

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/cache"
)

func RegisterRedisCache(id, prefix string, kvstroe *redis.Client) error {
	id = buildID(compCache, id)
	logger, err := Logger(compCache)
	if err != nil {
		return err
	}

	instance := cache.NewRedisCache(kvstroe, logger, prefix)
	comp := app.MakeComponent(instance, nil)
	App().RegisterComponent(comp, id)
	return nil
}

func RegisterMemoryCache(id string, size int) {
	instance := cache.NewMemoryCache(size)
	id = buildID(compCache, id)
	comp := app.MakeComponent(instance, nil)
	App().RegisterComponent(comp, id)
}

func Cache(id string) cache.Cache {
	id = buildID(compCache, id)
	instance, ok := App().Component(id)
	if !ok {
		panic(id + noRegister)
	}

	return instance.Object().(cache.Cache)
}
