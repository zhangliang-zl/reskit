package facade

import (
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/httpserver"
)

func RegisterServer(options httpserver.Options, id string) error {
	logger := Logger(compServer)
	engine, err := httpserver.New(options, logger, reskit.App().Context())
	if err != nil {
		return err
	}

	instance := component.Make(engine, engine.Run, engine.Close)
	return reskit.App().SetComponent(compServer, id, instance)
}

func Server(id string) *httpserver.Engine {
	res, ok := reskit.App().Component(compServer, id)

	if !ok {
		panic(compServer + noRegister)
	}

	return res.(*httpserver.Engine)
}
