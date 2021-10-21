package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit/redis"
	"time"
)

func main() {
	logger := log.NewHelper(log.DefaultLogger, log.WithMessageKey("redis"))
	rds, _ := redis.New(
		redis.WithLogger(logger),
		redis.WithPassword(""),
		redis.WithDbNum(2),
	)
	ctx := context.Background()
	rds.Set(ctx, "a", "1", time.Second*0)
}
