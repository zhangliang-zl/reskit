package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/logs"
	"time"
)

type Options struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	LogLevel string `json:"logLevel"`
}

func New(opts Options, logger logs.Logger) (*redis.Client, error) {
	if opts.LogLevel != "" {
		logger.SetLevel(logs.LevelVal(opts.LogLevel))
	}

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
