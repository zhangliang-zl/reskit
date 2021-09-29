package main

import (
	"context"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/exmple/grpcsrv/pb"
	"github.com/zhangliang-zl/reskit/facade"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/sd"
	"github.com/zhangliang-zl/reskit/service/server/grpcx"
	"log"
	"time"
)

type helloServer struct{}

func (s *helloServer) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	logs.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	app, err := reskit.Init(app.Options{
		Env:      "dev",
		LogLevel: logs.LevelDebug,
		Name:     "hero",
	})
	if err != nil {
		logs.Fatal(app.Run())
	}


	rpcServer := grpcx.NewServer()
	pb.RegisterHelloServer(rpcServer, &helloServer{})

	opts := grpcx.ServerOptions{
		Service: sd.Service{
			Name:                "hero",
			Tag:                 []string{"hero"},
			IP:                  "127.0.0.1",
			Port:                8086,
			HealthCheckInterval: time.Duration(10) * time.Second,
		},
		ConsulAddress: "127.0.0.1:8500",
	}

	if err := facade.RegisterGrpcServer(opts, "hero", rpcServer); err != nil {
		logs.Fatalln(err)
	}

	logs.Fatal(app.Run())
}
