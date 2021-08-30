package logs

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, msg string, data ...interface{})
	Info(ctx context.Context, msg string, data ...interface{})
	Warn(ctx context.Context, msg string, data ...interface{})
	Error(ctx context.Context, msg string, data ...interface{})
	Level() LogLevel
	SetLevel(level LogLevel)
}

type Factory interface {
	Get(tag string) (Logger, error)
}
