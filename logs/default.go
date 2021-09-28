package logs

import "github.com/zhangliang-zl/reskit/logs/driver/stdout"

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
