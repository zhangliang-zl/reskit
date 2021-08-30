package redis

import "context"
import "github.com/zhangliang-zl/reskit/logs"

type logging struct {
	l logs.Logger
}

func (r *logging) Printf(ctx context.Context, msg string, data ...interface{}) {
	r.l.Info(ctx, msg, data)
}
