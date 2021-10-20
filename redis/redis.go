package redis

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"time"
)

type Options struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func New(opts Options, logger *log.Helper) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return client, err
	}

	redis.SetLogger(&logWriter{l: logger})
	client.AddHook(logHook{
		l: logger,
	})
	return client, nil
}
