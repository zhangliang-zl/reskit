package logs

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs/driver/stdout"
)

type Logger interface {
	Debug(ctx context.Context, msg string, data ...interface{}) error
	Info(ctx context.Context, msg string, data ...interface{}) error
	Warn(ctx context.Context, msg string, data ...interface{}) error
	Error(ctx context.Context, msg string, data ...interface{}) error
	Fatal(ctx context.Context, msg string, data ...interface{}) error
	Level() Level
	SetLevel(level Level)
}

type Factory interface {
	Get(tag string) (Logger, error)
}

var DefaultFactory = NewFactory(LevelInfo, stdout.Driver())

func SetDefaultFactory(factory Factory) {
	DefaultFactory = factory
}

func DefaultLogger(tag string) (Logger, error) {
	return DefaultFactory.Get(tag)
}

func MustDefaultLogger(tag string) Logger {
	logger, err := DefaultFactory.Get(tag)
	if err != nil {
		panic(err)
	}
	return logger
}
