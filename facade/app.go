package facade

import (
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/logs"
)

var appInstance *app.App

func App() *app.App {
	if appInstance == nil {
		panic("Please call the init method before this")
	}
	return appInstance
}

func NewApp(opts ...app.Option) (*app.App, error) {
	ins, err := app.New(opts...)
	if err == nil {
		appInstance = ins
	}
	return ins, err
}

func Logger(tag string) (logs.Logger, error) {
	return App().LoggerFactory().Get(tag)
}
