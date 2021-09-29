package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

var ErrWriterIsNil = errors.New("writer is nil")

type logger struct {
	level  Level
	writer io.Writer
	sync.Mutex
}

func (l *logger) Debug(ctx context.Context, msg string, data ...interface{}) error {
	if l.level <= LevelDebug {
		return l.write(ctx, LevelDebug, msg, data...)
	}

	return nil
}

func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) error {
	if l.level <= LevelInfo {
		return l.write(ctx, LevelInfo, msg, data...)
	}
	return nil
}

func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) error {
	if l.level <= LevelWarn {
		return l.write(ctx, LevelWarn, msg, data...)
	}
	return nil
}

func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) error {
	if l.level <= LevelError {
		return l.write(ctx, LevelError, msg, data...)
	}
	return nil
}

func (l *logger) Fatal(ctx context.Context, msg string, data ...interface{}) error {
	if l.level <= LevelFatal {
		_ = l.write(ctx, LevelFatal, msg, data...)
		os.Exit(1)
	}
	return nil
}

func (l *logger) Level() Level {
	return l.level
}

func (l *logger) SetLevel(level Level) {
	l.level = level
}

func (l *logger) write(ctx context.Context, level Level, msg string, data ...interface{}) error {
	if l.writer == nil {
		return ErrWriterIsNil
	}

	traceID := GetTraceID(ctx)
	prefix := "[" + level.String() + "] tid[" + traceID + "] "

	if len(data) > 0 {
		msg = fmt.Sprintf(msg, data...)
	}

	_, err := l.writer.Write([]byte(prefix + msg + "\n"))
	return err
}

func NewLogger(level Level, writer io.Writer) Logger {
	return &logger{
		level:  level,
		writer: writer,
	}
}

