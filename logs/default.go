package logs

import (
	"log"
	"log/syslog"
)

func DefaultLogger(tag string) Logger {
	logger, err := DefaultLoggerFactory.Get(tag)
	if err != nil {
		log.Printf("logger load error:%s\n; auto switch stdout logger factory", err.Error())
		SwitchStdout(DefaultLevel)
	}

	return logger
}

var DefaultLevel = LevelInfo
var DefaultLoggerFactory = NewStdoutFactory(DefaultLevel)

func SwitchSyslog(level Level, pri syslog.Priority, projectName string) {
	DefaultLoggerFactory = NewSyslogFactory(level, pri, projectName)
}

func SwitchStdout(level Level) {
	DefaultLoggerFactory = NewStdoutFactory(level)
}
