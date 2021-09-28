package facade

import (
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/httpserver"
)

func RegisterHttpServer(options httpserver.Options, id string) error {
	logger, err := Logger(serviceHttpServer)
	if err != nil {
		return err
	}

	engine, err := httpserver.New(options, logger, App().Context())
	if err != nil {
		return err
	}

	instance := app.MakeService(engine, engine.Start, engine.Stop)
	id = buildID(serviceHttpServer, id)
	App().RegisterService(instance, id)
	return nil
}

func Server(id string) *httpserver.Engine {
	id = buildID(serviceHttpServer, id)
	res, ok := App().Service(id)

	if !ok {
		panic(serviceHttpServer + noRegister)
	}

	return res.Object().(*httpserver.Engine)
}
