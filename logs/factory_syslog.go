package logs

import (
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

func (factory *SyslogFactory) Get(tag string) (Logger, error) {
	v, ok := factory.Load(tag)
	if !ok {
		record, err := syslogrecord.NewRecorder(tag, syslogrecord.Priority(factory.priority))
		logger := NewLogger(record, factory.level, tag)
		if err != nil {
			factory.Store(tag, logger)
		}

		return logger, err
	}

	return v.(Logger), nil
}

func NewSyslogFactory(level Level, priority syslog.Priority, projectName string) LoggerFactory {
	return &SyslogFactory{
		level:       level,
		priority:    priority,
		projectName: projectName,
	}
}
