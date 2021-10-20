package redis

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"time"
)

type logHook struct {
	l *log.Helper
}

func (logHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx = context.WithValue(ctx, "redis-once-time", time.Now().UnixNano()/1e3)
	return ctx, nil
}

func (h logHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	start := ctx.Value("redis-once-time").(int64)
	elapsed := float64(time.Now().UnixNano()/1e3-start) / 1000
	h.l.Infof("%s cost time: %.3fms ", cmd.String(), elapsed)
	return nil
}

func (logHook) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (logHook) AfterProcessPipeline(_ context.Context, _ []redis.Cmder) error {
	return nil
}
