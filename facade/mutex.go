package facade

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/lock"
)

func RegisterMutexFactory(id string, kvstroe *redis.Client) error {
	id = buildID(compMutex, id)
	logger, err := Logger(compMutex)
	if err != nil {
		return err
	}
	factory := lock.NewRedisMutexFactory(logger, kvstroe, "Mutex:"+id)

	instance := app.MakeComponent(factory, nil)
	App().RegisterComponent(instance, id)
	return nil
}

func MutexFactory(id string) lock.Factory {
	id = buildID(compMutex, id)
	instance, ok := App().Component(id)

	if !ok {
		panic(compMutex + noRegister)
	}
	return instance.Object().(lock.Factory)
}
