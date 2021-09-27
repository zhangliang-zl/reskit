package facade

import (
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/grpcx"
	"google.golang.org/grpc"
)

func RegisterGrpcServer(opts grpcx.ServerOptions, id string, rpcServer *grpc.Server) error {
	logger := Logger(compGrpc)
	engine := grpcx.NewServer(opts, rpcServer, logger)
	instance := component.Make(engine, engine.Run, engine.Close)
	return reskit.App().SetComponent(compGrpc, id, instance)
}

func GrpcServer(id string) *grpcx.Server {
	res, ok := reskit.App().Component(compGrpc, id)

	if !ok {
		panic(compGrpc + noRegister)
	}

	return res.(*grpcx.Server)
}
