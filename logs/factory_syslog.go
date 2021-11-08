package logs

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs/stdout"
	syslogrecord "github.com/zhangliang-zl/reskit/logs/syslog"
	"log/syslog"
	"sync"
)

type SyslogFactory struct {
	level       Level
	priority    syslog.Priority
	projectName string
	sync.Map
}

func (factory *SyslogFactory) Get(tag string) Logger {
	v, ok := factory.Load(tag)
	if !ok {
		tag = factory.projectName + "/" + tag
		record, err := syslogrecord.NewRecorder(tag, syslogrecord.Priority(factory.priority))
		logger := NewLogger(record, factory.level, tag)
		if err != nil {
			// syslog 出错，使用stdout logger
			logger = NewLogger(stdout.NewRecorder(), DefaultLevel, tag)
			logger.Error(context.Background(), "syslog create error %s", err.Error())
		}

		factory.Store(tag, logger)
		return logger
	}

	return v.(Logger)
}

func NewSyslogFactory(level Level, priority syslog.Priority, projectName string) LoggerFactory {
	return &SyslogFactory{
		level:       level,
		priority:    priority,
		projectName: projectName,
	}
}
