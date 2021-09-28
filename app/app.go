package app

import (
	"context"
	"errors"
	"github.com/zhangliang-zl/reskit/logs"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	opts   options
	logger logs.Logger

	ctx    context.Context
	cancel func()

	components map[string]Component
	services   map[string]Service
}

func (app *App) Run() error {
	if len(app.services) == 0 {
		return nil
	}

	app.logger.Info(app.ctx, "start running ")
	eg, ctx := errgroup.WithContext(app.ctx)
	wg := sync.WaitGroup{}
	for _, srv := range app.services {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}

	wg.Wait()

	c := make(chan os.Signal, 1)
	signal.Notify(c, app.opts.signals...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				app.Stop()
			}
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		for _, comp := range app.components {
			if err := comp.Close(); err != nil {
				return err
			}
		}

		return err
	}
	return nil
}

func (app *App) Stop() {
	if app.cancel != nil {
		app.cancel()
	}
}

func (app *App) RegisterComponent(component Component, id string) {
	app.components[id] = component
}

func (app *App) RegisterService(service Service, id string) {
	app.services[id] = service
}

func (app *App) Component(id string) (Component, bool) {
	comp, ok := app.components[id]
	return comp, ok
}

func (app *App) Service(id string) (Service, bool) {
	svc, ok := app.services[id]
	return svc, ok
}

func (app *App) Version() string {
	return app.opts.env
}

func (app *App) Env() string {
	return app.opts.env
}

func (app *App) Name() string {
	return app.opts.name
}

func (app *App) LoggerFactory() logs.Factory {
	return app.opts.loggerFactory
}

func (app *App) Context() context.Context {
	return app.ctx
}

func New(opts ...Option) (*App, error) {

	o := options{
		loggerFactory: logs.DefaultFactory,
		ctx:           context.Background(),
		signals:       []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}

	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	logger, err := o.loggerFactory.Get("_app")

	app := &App{
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		opts:       o,
		services:   make(map[string]Service, 0),
		components: make(map[string]Component, 0),
	}

	if err != nil {
		return app, err
	}

	return app, nil
}
