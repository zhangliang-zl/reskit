package facade

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	kvstore "github.com/zhangliang-zl/reskit/redis"
)

func RegisterRedis(opts kvstore.Options, id string) error {
	logger := Logger(compRedis)
	client, err := kvstore.New(opts, logger)
	if err != nil {
		return err
	}

	instance := component.Make(client, nil, client.Close)
	return reskit.App().SetComponent(compRedis, id, instance)
}

func Redis(id string) *redis.Client {
	res, ok := reskit.App().Component(compRedis, id)
	if !ok {
		panic(compRedis + noRegister)
	}
	return res.(*redis.Client)
}
