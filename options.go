package reskit

import (
	"context"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/service"
	"os"
)

type Option func(o *options)

type options struct {
	env     string
	name    string
	version string

	components map[string]component.Interface
	services   map[string]service.Interface

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

func WithComponent(id string, c component.Interface) Option {
	return func(o *options) {
		if c != nil && id != "" {
			if o.components == nil {
				o.components = make(map[string]component.Interface)
			}
			o.components[id] = c
		}
	}
}

func WithService(id string, s service.Interface) Option {
	return func(o *options) {
		if s != nil && id != "" {
			if o.services == nil {
				o.services = make(map[string]service.Interface)
			}
			o.services[id] = s
		}
	}
}
