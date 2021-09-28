package main

import (
	"context"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/facade"
	"github.com/zhangliang-zl/reskit/grpcx"
	"github.com/zhangliang-zl/reskit/grpcx/sd"
	"github.com/zhangliang-zl/reskit/grpcx/test/pb"
	"github.com/zhangliang-zl/reskit/logs"
	"google.golang.org/grpc"
	"log"
	"time"
)

type helloServer struct{}

func (s *helloServer) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	app, err := reskit.Init(app.Options{
		Env:      "dev",
		LogLevel: logs.LevelDebug,
		Name:     "hero",
	})
	if err != nil {
		log.Fatal(app.Run())
	}


	rpcServer := grpc.NewServer()
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
		log.Fatalln(err)
	}

	log.Fatal(app.Run())
}
