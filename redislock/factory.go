package redislock

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/logs"
)

var DefaultLogger = logs.DefaultLogger("lock")

type FactoryOption func(factory *Factory)

func Logger(logger logs.Logger) FactoryOption {
	return func(factory *Factory) {
		factory.logger = logger
	}
}

func NewFactory(rdsClient *redis.Client, opts ...FactoryOption) *Factory {
	f := &Factory{
		logger:    DefaultLogger,
		rdsClient: rdsClient,
	}
	for _, opt := range opts {
		opt(f)
	}

	return f
}

type Factory struct {
	rdsClient *redis.Client
	logger    logs.Logger
}

func (f *Factory) New(key string, opts ...Option) *Mutex {
	var o = &Options{
		duration:      DefaultLocked,
		lockWaiting:   DefaultLockWaiting,
		retryInterval: DefaultRetryInterval,
		keyPrefix:     DefaultKeyPrefix,
	}

	for _, opt := range opts {
		opt(o)
	}

	return &Mutex{
		redisClient: f.rdsClient,
		logger:      f.logger,
		opts:        o,
		key:         o.keyPrefix + key,
	}
}
