package service

import (
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/transport/grpcx"
	"google.golang.org/grpc"
)

func NewGrpcServer(options grpcx.ServerOptions, grpcServer *grpc.Server, logger logs.Logger) Interface {
	engine := grpcx.NewServer(options, grpcServer, logger)
	return Make(engine, engine.Start, engine.Stop)
}