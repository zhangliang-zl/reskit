package reskit

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	opts   *Options
	ctx    context.Context
	cancel context.CancelFunc
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func (a *App) Run() error {
	if len(a.opts.beforeStart) > 0 {
		for _, fn := range a.opts.beforeStart {
			if err := fn(); err != nil {
				return err
			}
		}
	}

	eg, ctx := errgroup.WithContext(a.ctx)
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
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

	if len(a.opts.afterStart) > 0 {
		for _, fn := range a.opts.afterStart {
			if err := fn(); err != nil {
				return err
			}
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				err := a.stop()
				if err != nil {
					a.opts.logger.Errorf("failed to stop app: %v", err)
					return err
				}
			}
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	if len(a.opts.afterStop) > 0 {
		for _, fn := range a.opts.afterStop {
			if err := fn(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *App) Name() string {
	return a.opts.name
}

func (a *App) Env() string {
	return a.opts.env
}

func (a *App) stop() error {
	a.cancel()

	for _, fn := range a.opts.beforeStop {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func New(opts ...Option) *App {
	o := &Options{
		env:     "dev",
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		logger:  log.NewHelper(log.With(log.DefaultLogger, "tag", "app")),
		servers: make([]Server, 0),
		ctx:     context.Background(),
	}
	for _, opt := range opts {
		opt(o)
	}
	ctx, cancel := context.WithCancel(o.ctx)

	return &App{
		opts:   o,
		cancel: cancel,
		ctx:    ctx,
	}
}
