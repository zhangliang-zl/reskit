package logs

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs/empty"
	"github.com/zhangliang-zl/reskit/logs/stdout"
	"github.com/zhangliang-zl/reskit/logs/syslog"
)

type Recorder interface {
	Record(m string)
}

type Logger interface {
	Debug(ctx context.Context, msg string, data ...interface{})
	Info(ctx context.Context, msg string, data ...interface{})
	Warn(ctx context.Context, msg string, data ...interface{})
	Error(ctx context.Context, msg string, data ...interface{})
	Fatal(ctx context.Context, msg string, data ...interface{})
}

type LoggerFactory interface {
	Get(tag string) (Logger, error)
}

var DefaultRecorder = stdout.NewRecorder()
var DefaultLevel = LevelInfo

func DefaultLogger(tag string) Logger {
	return NewLogger(DefaultRecorder, DefaultLevel, tag)
}

var _ Recorder = (*stdout.Recorder)(nil)
var _ Recorder = (*syslog.Recorder)(nil)
var _ Recorder = (*empty.Recorder)(nil)