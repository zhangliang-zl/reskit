package logs

import (
	"context"
	"fmt"
	"github.com/zhangliang-zl/reskit/logs/driver"
	"sync"
)

type logger struct {
	level  LogLevel
	writer driver.Writer
	sync.Mutex
}

func (d *logger) Debug(ctx context.Context, msg string, data ...interface{}) {
	if d.level >= DEBUG {
		d.record(ctx, DEBUG, msg, data...)
	}
}

func (d *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if d.level >= INFO {
		d.record(ctx, INFO, msg, data...)
	}
}

func (d *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if d.level >= WARN {
		d.record(ctx, WARN, msg, data...)
	}
}

func (d *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	d.record(ctx, ERROR, msg, data...)
}

func (d *logger) Level() LogLevel {
	return d.level
}

func (d *logger) SetLevel(level LogLevel) {
	d.level = level
}

func (d *logger) record(ctx context.Context, level LogLevel, msg string, data ...interface{}) {
	if d.writer == nil {
		return
	}

	traceID := findTraceID(ctx)
	prefix := "[" + LevelName(level) + "] tid[" + traceID + "] "

	if len(data) > 0 {
		msg = fmt.Sprintf(msg, data...)
	}

	d.writer.Record(prefix + msg)
}

func NewLogger(level LogLevel, recorder driver.Writer) Logger {
	return &logger{
		level:  level,
		writer: recorder,
	}
}
