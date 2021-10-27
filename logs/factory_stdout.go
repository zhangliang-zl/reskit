package logs

import (
	"github.com/zhangliang-zl/reskit/logs/stdout"
	"log/syslog"
	"sync"
)

type StdoutFactory struct {
	level       Level
	priority    syslog.Priority
	projectName string
	sync.Map
}

func (factory *StdoutFactory) Get(tag string) (Logger, error) {
	v, ok := factory.Load(tag)
	if !ok {

		logger := NewLogger(stdout.NewRecorder(), factory.level, tag)
		factory.Store(tag, logger)
		return logger, nil
	}

	return v.(Logger), nil
}

func NewStdoutFactory(level Level ) LoggerFactory {
	return &StdoutFactory{
		level:       level,
	}
}
