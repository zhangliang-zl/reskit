package syslog

import (
	"log/syslog"
)

type Recorder struct {
	w *syslog.Writer
	p syslog.Priority
}

func (r *Recorder) Record(m string) {
	_ = r.w.Info(m)
}

type Option func(*Recorder)

var DefaultPriority = syslog.LOG_DEBUG | syslog.LOG_LOCAL6

func Priority(priority syslog.Priority) Option {
	return func(r *Recorder) {
		r.p = priority
	}
}

func NewRecorder(tag string, opts ...Option) (*Recorder, error) {
	r := &Recorder{
		p: DefaultPriority,
	}

	for _, opt := range opts {
		opt(r)
	}

	w, err := syslog.New(r.p, tag)
	if err == nil {
		r.w = w
	}

	return r, err
}
