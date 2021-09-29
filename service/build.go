package service

import (
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/service/server/grpcx"
	"github.com/zhangliang-zl/reskit/service/server/web"
	"google.golang.org/grpc"
)

func NewHttpServer(options web.Options, logger logs.Logger) Interface {
	engine := web.NewServer(options, logger)
	return Make(engine, engine.Start, engine.Stop)
}

func NewGrpcServer(options grpcx.ServerOptions, grpcServer *grpc.Server, logger logs.Logger) Interface {
	engine := grpcx.NewServer(options, grpcServer, logger)
	return Make(engine, engine.Start, engine.Stop)
}
