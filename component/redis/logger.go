package redis

import "context"
import "github.com/zhangliang-zl/reskit/log"

type logWriter struct {
	l log.Logger
}

func (r *logWriter) Printf(ctx context.Context, msg string, data ...interface{}) {
	r.l.Info(ctx, msg, data)
}
