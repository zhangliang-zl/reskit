package facade

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/lock"
)

func RegisterMutexFactory(id string, kvstroe *redis.Client) error {
	logger := Logger(compMutex)
	factory := lock.NewFactoryForRedis(logger, kvstroe)
	instance := component.Make(factory, nil, nil)
	return reskit.App().SetComponent(compMutex, id, instance)
}

func MutexFactory(id string) lock.Factory {
	instance, ok := reskit.App().Component(compMutex, id)

	if !ok {
		panic(compMutex + noRegister)
	}
	return instance.(lock.Factory)
}
