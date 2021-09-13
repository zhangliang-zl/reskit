package logs

import (
	"strings"
)

type Level uint8

const (
	LevelError Level = 0
	LevelWarn  Level = 1
	LevelInfo  Level = 2
	LevelDebug Level = 3
)

var levelNamesMapping = map[Level]string{
	LevelError: "error",
	LevelWarn:  "warn",
	LevelInfo:  "info",
	LevelDebug: "debug",
}

func LevelVal(label string) Level {
	label = strings.ToLower(label)
	for l, name := range levelNamesMapping {
		if name == label {
			return l
		}
	}

	return LevelError
}

func LevelName(l Level) string {
	name, ok := levelNamesMapping[l]
	if ok {
		return name
	}

	return "error"
}
