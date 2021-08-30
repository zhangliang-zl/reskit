package logs

import "context"

type EmptyLogger struct {
}

func (e EmptyLogger) Debug(ctx context.Context, msg string, data ...interface{}) {
}

func (e EmptyLogger) Info(ctx context.Context, msg string, data ...interface{}) {
}

func (e EmptyLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
}

func (e EmptyLogger) Error(ctx context.Context, msg string, data ...interface{}) {
}

func (e EmptyLogger) Level() LogLevel {
	return DEBUG
}

func (e EmptyLogger) SetLevel(level LogLevel) {
}
