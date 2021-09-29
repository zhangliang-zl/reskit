package grpcx

import (
	"context"
	"github.com/zhangliang-zl/reskit/sd/consul"
	"google.golang.org/grpc"
	"time"
)

type ClientOptions struct {
	ConsulAddress string
	Target        string
	Timeout       time.Duration
}

func (t ClientOptions) ConsulTarget() string {
	return "consul://" + t.ConsulAddress + "/" + t.Target
}

func Dial(ctx context.Context, opts ClientOptions, dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, _ = context.WithTimeout(ctx, opts.Timeout)
	sdIns := consul.New(opts.ConsulAddress)

	dialOpts = append(
		dialOpts, grpc.WithResolvers(sdIns.ResolverBuilder()),
	)

	conn, err := grpc.DialContext(ctx, opts.ConsulTarget(), dialOpts...)
	return conn, err
}
