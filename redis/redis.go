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
	Logger   *log.Helper
}

var (
	DefaultAddress  = "127.0.0.1:6379"
	DefaultDB       = 0
	DefaultPassword = ""
	DefaultLogger   = log.NewHelper(log.DefaultLogger, log.WithMessageKey("redis"))
)

type Option func(options *Options)

func New(opts ...Option) (*redis.Client, error) {
	o := &Options{
		Addr:     DefaultAddress,
		DB:       DefaultDB,
		Password: DefaultPassword,
		Logger:   DefaultLogger,
	}

	for _, opt := range opts {
		opt(o)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     o.Addr,
		Password: o.Password,
		DB:       o.DB,
	})

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return client, err
	}

	redis.SetLogger(&logWriter{l: o.Logger})
	client.AddHook(logHook{
		l: o.Logger,
	})
	return client, nil
}
