package app

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs"
	"os"
)

type Option func(o *options)

type options struct {
	env     string
	name    string
	version string

	ctx           context.Context
	signals       []os.Signal
	loggerFactory logs.Factory
}

func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

func WithVersion(version string) Option {
	return func(o *options) { o.version = version }
}

func WithEnv(env string) Option {
	return func(o *options) {
		o.env = env
	}
}

func WithLoggerFactory(factory logs.Factory) Option {
	return func(o *options) {
		if factory != nil {
			o.loggerFactory = factory
		}
	}
}

func WithSignal(signals ...os.Signal) Option {
	return func(o *options) {
		if o.signals != nil {
			o.signals = signals
		}
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		if ctx != nil {
			o.ctx = ctx
		}
	}
}
