package logs

import (
	"log/syslog"
)

func DefaultLogger(tag string) Logger {
	return DefaultLoggerFactory.Get(tag)
}

var DefaultLevel = LevelInfo

var DefaultLoggerFactory = NewStdoutFactory(DefaultLevel)

func SwitchSyslogFactory(l Level, pri syslog.Priority, projectName string) {
	DefaultMode = ModeSyslog
	DefaultLoggerFactory = NewSyslogFactory(l, pri, projectName)
}

func SwitchStdoutFactory(l Level) {
	DefaultMode = ModeStdout
	DefaultLoggerFactory = NewStdoutFactory(l)
}

const (
	ModeStdout = "stdout"
	ModeSyslog = "syslog"
)

var DefaultMode = ModeStdout
