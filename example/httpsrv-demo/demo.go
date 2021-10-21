package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit/redis"
	"time"
)

func main() {
	logger := log.NewHelper(log.With(log.DefaultLogger, "tag", "redis"))
	rds, _ := redis.New(
		redis.Logger(logger),
		redis.Password(""),
		redis.DbNum(2),
	)
	ctx := context.Background()
	rds.Set(ctx, "a", "1", time.Second*0)
}
