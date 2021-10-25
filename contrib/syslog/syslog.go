package syslog

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"log/syslog"
	"sync"
)

type logger struct {
	w    *syslog.Writer
	pool *sync.Pool
}

func (l *logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	buf.WriteString(level.String())
	for i := 0; i < len(keyvals); i += 2 {
		_, _ = fmt.Fprintf(buf, " %s=%v", keyvals[i], keyvals[i+1])
	}
	_ = l.w.Info(buf.String()) //nolint:gomnd
	l.pool.Put(buf)
	return nil
}

func NewLogger(priority syslog.Priority, tag string) (log.Logger, error) {
	w, err := syslog.New(priority, tag)
	if err != nil {
		return nil, err
	}

	return &logger{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		w: w,
	}, nil
}
