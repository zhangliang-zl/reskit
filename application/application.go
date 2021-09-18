package application

import (
	"context"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/logs"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Application struct {
	env      string
	pid      string // Save pid to this file
	logLevel string
	prjName  string

	components component.Container

	logger logs.Logger
	ctx    context.Context
}

func (app *Application) Init(options Options) error {

	app.env = options.Env
	app.pid = options.PID
	app.logLevel = options.LogLevel
	app.prjName = options.PrjName

	app.logger, _ = loggerFactory.Get("_application")
	app.ctx = context.Background()
	app.components = component.NewContainer(app.logger, app.ctx)

	return nil
}

func (app *Application) Run() error {
	if err := app.savePID(app.pid); err != nil {
		return err
	}

	app.logger.Info(app.ctx, "start running ")

	if err := app.components.Run(); err != nil {
		app.logger.Info(app.ctx, "components run err %s", err.Error())
		return err
	}

	// Waiting close signal.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan

	if err := app.Close(); err != nil {
		app.logger.Error(app.ctx, "close error %s", err.Error())
	} else {
		app.logger.Info(app.ctx, "close finish.")
	}

	return nil
}

func (app *Application) Close() error {
	return app.components.Close()
}

func (app *Application) savePID(file string) error {
	if file == "" {
		return nil
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(strconv.Itoa(os.Getpid()))
	return err
}

func (app *Application) Component(cType string, id string) (interface{}, bool) {
	key := component.BuildKey(cType, id)
	return app.components.Get(key)
}

func (app *Application) SetComponent(cType string, id string, val component.Interface) error {
	key := component.BuildKey(cType, id)
	return app.components.Set(key, val)
}

func (app *Application) Env() string {
	return app.env
}

func (app *Application) LoggerLevel() string {
	return app.logLevel
}

func (app *Application) PrjName() string {
	return app.prjName
}

func (app *Application) Context() context.Context {
	return app.ctx
}

func NewApplication(options Options) (*Application, error) {
	app := &Application{}
	err := app.Init(options)
	return app, err
}
