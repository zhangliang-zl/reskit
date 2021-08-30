package logs

import (
	"strings"
)

type LogLevel uint8

const (
	ERROR LogLevel = 0
	WARN  LogLevel = 1
	INFO  LogLevel = 2
	DEBUG LogLevel = 3
)

var levelLabels = map[LogLevel]string{
	ERROR: "error",
	WARN:  "warn",
	INFO:  "info",
	DEBUG: "debug",
}

func LevelVal(label string) LogLevel {
	label = strings.ToLower(label)
	for l, name := range levelLabels {
		if name == label {
			return l
		}
	}

	return ERROR
}

func LevelName(l LogLevel) string {
	name, ok := levelLabels[l]
	if ok {
		return name
	}

	return "error"
}
