package redis

import "context"
import "github.com/zhangliang-zl/reskit/logs"

type logWriter struct {
	l logs.Logger
}

func (r *logWriter) Printf(ctx context.Context, msg string, data ...interface{}) {
	r.l.Info(ctx, msg, data)
}
