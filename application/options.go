package application

import "github.com/zhangliang-zl/reskit/logs"

type Options struct {
	Env, PrjName, PID string
	LogLevel          logs.Level
}
