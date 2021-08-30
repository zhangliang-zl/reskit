package reskit

import (
	"github.com/zhangliang-zl/reskit/application"
)

var app *application.Application
var isInit bool

func Init(options application.Options) (*application.Application, error) {
	var err error
	app, err = application.NewApplication(options)
	if err == nil {
		isInit = true
	}

	return app, err
}

func App() *application.Application {
	if !isInit {
		panic("Please call the init method before this")
	}
	return app
}
