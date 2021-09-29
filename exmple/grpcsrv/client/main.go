package main

import (
	"context"
	"github.com/zhangliang-zl/reskit/exmple/grpcsrv/pb"
	"github.com/zhangliang-zl/reskit/service/server/grpcx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"log"
	"os"
	"time"
)

func main() {

	ctx := context.Background()
	opts := grpcx.ClientOptions{
		ConsulAddress: "127.0.0.1:8500",
		Target:        "hero",
		Timeout:       30 * time.Second,
	}

	conn, err := grpcx.Dial(ctx, opts, grpcx.WithBlock(), grpcx.WithInsecure(), grpcx.WithBalancerName(roundrobin.Name))
	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewHelloClient(conn)

	// Contact the server1 and print out its response.
	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.Message)
		time.Sleep(time.Second * 1)
	}
}
