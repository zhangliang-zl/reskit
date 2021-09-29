package service

import (
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/transport/grpcx"
	"github.com/zhangliang-zl/reskit/transport/httpx"
	"google.golang.org/grpc"
)

func NewHttpServer(options httpx.Options, logger logs.Logger) Interface {
	engine := httpx.NewServer(options, logger)
	return Make(engine, engine.Start, engine.Stop)
}

func NewGrpcServer(options grpcx.ServerOptions, grpcServer *grpc.Server, logger logs.Logger) Interface {
	engine := grpcx.NewServer(options, grpcServer, logger)
	return Make(engine, engine.Start, engine.Stop)
}
