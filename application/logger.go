package application

import (
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/logs/driver/stdout"
)

var loggerFactory = logs.NewFactory("debug", stdout.Driver())

func SetLoggerFactory(factory logs.Factory) {
	loggerFactory = factory
}

func Logger(tag string) (logs.Logger, error) {
	return loggerFactory.Get(tag)
}
