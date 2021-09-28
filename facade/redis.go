package facade

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/app"
	kvstore "github.com/zhangliang-zl/reskit/redis"
)

func RegisterRedis(opts kvstore.Options, id string) error {
	id = buildID(compRedis, id)
	logger, err := Logger(compRedis)
	if err != nil {
		return err
	}

	instance, err := kvstore.New(opts, logger)
	if err != nil {
		return err
	}

	comp := app.MakeComponent(instance, instance.Close)
	appInstance.RegisterComponent(comp, id)

	return nil
}

func Redis(id string) *redis.Client {
	id = buildID(compRedis, id)
	res, ok := App().Component(id)
	if !ok {
		panic(compRedis + noRegister)
	}
	return res.Object().(*redis.Client)
}
