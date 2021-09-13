package logs

import (
	"context"
	"fmt"
	"github.com/zhangliang-zl/reskit/logs/driver"
	"sync"
)

type Logger interface {
	Debug(ctx context.Context, msg string, data ...interface{})
	Info(ctx context.Context, msg string, data ...interface{})
	Warn(ctx context.Context, msg string, data ...interface{})
	Error(ctx context.Context, msg string, data ...interface{})
	Level() Level
	SetLevel(level Level)
}

type logger struct {
	level  Level
	writer driver.Writer
	sync.Mutex
}

func (d *logger) Debug(ctx context.Context, msg string, data ...interface{}) {
	if d.level >= LevelDebug {
		d.record(ctx, LevelDebug, msg, data...)
	}
}

func (d *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if d.level >= LevelInfo {
		d.record(ctx, LevelInfo, msg, data...)
	}
}

func (d *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if d.level >= LevelWarn {
		d.record(ctx, LevelWarn, msg, data...)
	}
}

func (d *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	d.record(ctx, LevelError, msg, data...)
}

func (d *logger) Level() Level {
	return d.level
}

func (d *logger) SetLevel(level Level) {
	d.level = level
}

func (d *logger) record(ctx context.Context, level Level, msg string, data ...interface{}) {
	if d.writer == nil {
		return
	}

	traceID := GetTraceID(ctx)
	prefix := "[" + LevelName(level) + "] tid[" + traceID + "] "

	if len(data) > 0 {
		msg = fmt.Sprintf(msg, data...)
	}

	d.writer.Record(prefix + msg)
}

func NewLogger(level Level, writer driver.Writer) Logger {
	return &logger{
		level:  level,
		writer: writer,
	}
}
