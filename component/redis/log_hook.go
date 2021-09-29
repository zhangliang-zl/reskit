package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zhangliang-zl/reskit/log"
	"time"
)

type logHook struct {
	l log.Logger
}

func (logHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx = context.WithValue(ctx, "redis-once-time", time.Now().UnixNano()/1e3)
	return ctx, nil
}

func (h logHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	start := ctx.Value("redis-once-time").(int64)
	elapsed := float64(time.Now().UnixNano()/1e3-start) / 1000
	h.l.Info(ctx, fmt.Sprintf("%s usetime: %.3fms ", cmd.String(), elapsed))
	return nil
}

func (logHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (logHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
