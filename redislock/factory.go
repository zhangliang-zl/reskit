package redislock

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

var DefaultLogger = log.NewHelper(log.DefaultLogger, log.WithMessageKey("lock"))

type FactoryOption func(factory *Factory)

func Logger(logger *log.Helper) FactoryOption {
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
	logger    *log.Helper
}

func (f *Factory) New(key string, opts ...Option) *Mutex {
	var o = &Options{
		Duration:      DefaultLocked,
		LockWaiting:   DefaultLockWaiting,
		RetryInterval: DefaultRetryInterval,
		KeyPrefix:     DefaultKeyPrefix,
	}

	for _, opt := range opts {
		opt(o)
	}

	return &Mutex{
		redisClient: f.rdsClient,
		logger:      f.logger,
		opts:        o,
		key:         o.KeyPrefix + key,
	}
}
