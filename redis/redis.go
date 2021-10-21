package redis

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"time"
)

type Options struct {
	addr     string
	password string
	dbNum    int
	logger   *log.Helper
}

func New(opts ...Option) (*redis.Client, error) {
	o := DefaultOption

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
