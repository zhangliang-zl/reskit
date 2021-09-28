package facade

import (
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/grpcx"
	"google.golang.org/grpc"
)

func RegisterGrpcServer(opts grpcx.ServerOptions, id string, rpcServer *grpc.Server) error {
	id = buildID(serviceGrpc, id)
	logger, err := Logger(serviceGrpc)
	if err != nil {
		return err
	}

	engine := grpcx.NewServer(opts, rpcServer, logger)
	instance := app.MakeService(engine, engine.Start, engine.Stop)
	App().RegisterService(instance, id)
	return nil
}

func GrpcServer(id string) *grpcx.Server {
	id = buildID(serviceGrpc, id)
	res, ok := App().Service(id)

	if !ok {
		panic(serviceGrpc + noRegister)
	}

	return res.Object().(*grpcx.Server)
}
