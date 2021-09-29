package service

import (
	"context"
	"fmt"
	"github.com/zhangliang-zl/reskit/component/snowflake"
	"github.com/zhangliang-zl/reskit/log"
	"net"
	"strconv"
	"time"
)

type Server struct {
	port     int
	logger   log.Logger
	idWorker *snowflake.Worker
}

type ServerOptions struct {
	Port       int // 推荐 5101
	WorkerID   int64
	WorkerBits int64 // 推荐 4~5
	NumberBits int64 // 推荐18
	Epoch      int64 // 推荐系统开始使用时开始
}

func NewServer(opts ServerOptions, logger log.Logger) (*Server, error) {
	workerID := opts.WorkerID
	idWorker, err := snowflake.NewWorker(workerID, opts.WorkerBits, opts.NumberBits, opts.Epoch)
	if err != nil {
		return nil, err
	}

	return &Server{
		port:     opts.Port,
		idWorker: idWorker,
		logger:   logger,
	}, nil
}

func (s *Server) Run() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	ctx := log.WithTraceID(context.Background())

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error(ctx, "Accept Error:%s", err.Error())
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) error {
	err := conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	if err != nil {
		return err
	}

	defer conn.Close()
	request := make([]byte, 128)
	ctx := log.WithTraceID(context.Background())
loop:
	for {
		readLen, err := conn.Read(request)
		if err != nil {
			continue
		}

		instruction := string(request[:readLen])
		if readLen == 0 {
			s.logger.Error(ctx, "ReadLen is 0 , client disconnect automatically")
			break
		}

		switch instruction {
		case "get uuid\r\n":
			uuid := s.idWorker.NextID()
			uuidStr := strconv.FormatInt(uuid, 10)
			respInfo := fmt.Sprintf("VALUE uuid 0 %d\r\n%s\r\nEND\r\n", len(uuidStr), uuidStr)
			conn.Write([]byte(respInfo))
			s.logger.Debug(ctx, "nextID is %d", uuid)

		case "quit\r\n":
			s.logger.Debug(ctx, "client is quit")
			break loop

		default:
			s.logger.Warn(ctx, "unknown instructions")
		}
		request = make([]byte, 128)
	}

	return nil
}
