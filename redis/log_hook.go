package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/logs"
	"time"
)

type logHook struct {
	l logs.Logger
}

func (logHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx = context.WithValue(ctx, "redis-once-time", time.Now().UnixNano()/1e3)
	return ctx, nil
}

func (h logHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	start := ctx.Value("redis-once-time").(int64)
	elapsed := float64(time.Now().UnixNano()/1e3-start) / 1000
	h.l.Info(ctx, "%s cost time: %.3fms ", cmd.String(), elapsed)
	return nil
}

func (logHook) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (logHook) AfterProcessPipeline(_ context.Context, _ []redis.Cmder) error {
	return nil
}
