package logs

import (
	"context"
	"fmt"
	"github.com/zhangliang-zl/reskit/logs/stdout"
	"os"
)

type logger struct {
	level Level
	w     Recorder
	tag   string
}

func (d *logger) Debug(ctx context.Context, msg string, data ...interface{}) {
	if d.level <= LevelDebug {
		d.record(ctx, LevelDebug, msg, data...)
	}
}

func (d *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if d.level <= LevelInfo {
		d.record(ctx, LevelInfo, msg, data...)
	}
}

func (d *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if d.level <= LevelWarn {
		d.record(ctx, LevelWarn, msg, data...)
	}
}

func (d *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if d.level <= LevelError {
		d.record(ctx, LevelError, msg, data...)
	}
}
func (d *logger) Fatal(ctx context.Context, msg string, data ...interface{}) {
	if d.level <= LevelFatal {
		d.record(ctx, LevelFatal, msg, data...)
		os.Exit(1)
	}
}

func (d *logger) record(ctx context.Context, level Level, msg string, data ...interface{}) {
	if d.w == nil {
		return
	}

	traceID := findTraceID(ctx)

	if len(data) > 0 {
		msg = fmt.Sprintf(msg, data...)
	}

	prefix := fmt.Sprintf("level=%s, tag=%s, trace_id=%s, msg=%s", level.String(), d.tag, traceID, msg)

	d.w.Record(prefix + msg)
}

func NewLogger(recorder Recorder, level Level, tag string) Logger {
	return &logger{
		level: level,
		w:     recorder,
		tag:   tag,
	}
}

var DefaultTag = ""
var DefaultLogger = NewLogger(stdout.NewRecorder(), LevelInfo, DefaultTag)
