package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/logs"
	"time"
)

type Options struct {
	addr     string
	password string
	dbNum    int
	logger   logs.Logger
}

func New(opts ...Option) (*redis.Client, error) {

	var o = &Options{
		addr:     "127.0.0.1:6379",
		password: "",
		dbNum:    0,
		logger:   logs.DefaultLogger("_redis"),
	}

	for _, opt := range opts {
		opt(o)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     o.addr,
		Password: o.password,
		DB:       o.dbNum,
	})

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return client, err
	}

	redis.SetLogger(&logWriter{l: o.logger})
	client.AddHook(logHook{
		l: o.logger,
	})
	return client, nil
}
