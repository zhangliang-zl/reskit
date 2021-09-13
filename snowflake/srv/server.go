package snowflake

import (
	"context"
	"fmt"
	"github.com/zhangliang-zl/reskit/logs"
	"net"
	"strconv"
	"time"
)

type Server struct {
	port     int
	logger   logs.Logger
	idWorker *Worker
}

type ServerOptions struct {
	Port       int // 推荐 5101
	MysqlDsn   string
	WorkerBits int64 // 推荐 4~5
	NumberBits int64 // 推荐18
	Epoch      int64 // 推荐系统开始使用时开始
}

func NewServer(opts ServerOptions, logger logs.Logger) (*Server, error) {
	workerFactory := NewWorkerIDFactory(opts.MysqlDsn)
	workerID, err := workerFactory.WorkID()
	if err != nil {
		return nil, err
	}

	idWorker, err := NewWorker(workerID, opts.WorkerBits, opts.NumberBits, opts.Epoch)
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
	ctx := logs.WithTraceID(context.Background())

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
	ctx := logs.WithTraceID(context.Background())
	for {
		readLen, err := conn.Read(request)
		if err != nil {
			s.logger.Error(ctx, "Read Error :%s", err.Error())
			continue
		}

		instruction := string(request[:readLen])
		if readLen == 0 {
			s.logger.Error(ctx, "ReadLen is 0 , client disconnect automatically")
			break
		} else if instruction == "uuid\r\n" {
			uuid := s.idWorker.NextID()
			uuidStr := strconv.FormatInt(uuid, 10) + "\r\n"
			conn.Write([]byte(uuidStr))
			s.logger.Debug(ctx, "nextID is "+uuidStr)
		} else if instruction == "quit\r\n" {
			s.logger.Debug(ctx, "client is quit")
			break
		} else {
			s.logger.Warn(ctx, "unknown instructions")
		}
		request = make([]byte, 128)
	}

	return nil
}
