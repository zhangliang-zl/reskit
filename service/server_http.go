package service

import (
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/transport/web"
)

func NewHttpServer(options web.Options, logger logs.Logger) Interface {
	engine := web.NewServer(options, logger)
	return Make(engine, engine.Start, engine.Stop)
}
