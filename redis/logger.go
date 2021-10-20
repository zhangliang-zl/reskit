package redis

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type logWriter struct {
	l *log.Helper
}

func (r *logWriter) Printf(_ context.Context, msg string, data ...interface{}) {
	r.l.Infof(msg, data)
}
