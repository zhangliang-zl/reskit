package grpcx

import (
	"context"
	"fmt"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/sd"
	"github.com/zhangliang-zl/reskit/sd/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"sync"
	"time"
)

type ServerOptions struct {
	Service       sd.Service
	ConsulAddress string
}

type Server struct {
	opts ServerOptions
	mu   sync.Mutex
	grpcServer *grpc.Server
	logger     logs.Logger
}

func NewServer(opts ServerOptions, grpcServer *grpc.Server, logger logs.Logger) *Server {
	return &Server{
		opts:       opts,
		grpcServer: grpcServer,
		logger:     logger,
	}
}

func (s *Server) GrpcServer() *grpc.Server {
	return s.grpcServer
}

func (s *Server) Start(ctx context.Context) error {
	logger := s.logger
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opts.Service.Port))
	if err != nil {
		logger.Error(ctx, "listen error: %v", err)
		return err
	}

	sdManager := consul.New(s.opts.ConsulAddress)
	if err := sdManager.Register(s.opts.Service); err != nil {
		logger.Error(ctx, "service discovery error: %v", err)
		return err
	}

	grpc_health_v1.RegisterHealthServer(s.grpcServer, &DefaultHealthServer{})

	afterCh := time.After(time.Second * 3)
	finalError := make(chan error, 1)

	go func(errorChan chan error) {
		var err error
		var retry = 5
		for i := 0; i < retry; i++ {
			if err = s.grpcServer.Serve(lis); err != nil {
				if err == grpc.ErrServerStopped {
					logger.Info(ctx, "grpcx server closed")
					break
				}

				logger.Warn(ctx, "grpcServer listenAdnServe. current %d,maxRetry %d ,err:%v", i+1, retry, err)
				time.Sleep(time.Millisecond * 500)
				continue
			}
		}

		errorChan <- err

	}(finalError)

	select {
	case <-afterCh:
		return nil
	case err := <-finalError:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.grpcServer.GracefulStop()
	s.mu.Unlock()
	return nil
}
