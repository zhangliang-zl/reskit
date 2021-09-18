package facade

import (
	"github.com/zhangliang-zl/reskit/application"
	"github.com/zhangliang-zl/reskit/logs"
)

func Logger(tag string) logs.Logger {
	l, err := application.Logger(tag)
	if err != nil {
		panic(err)
	}
	return l
}
