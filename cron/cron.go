package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/zhangliang-zl/reskit"
	"time"
)

type Server struct {
	Cron         *cron.Cron
	closeTimeout time.Duration
}

func (s *Server) Start(_ context.Context) error {
	s.Cron.Run()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	ctx := s.Cron.Stop()
	timer := time.NewTimer(s.closeTimeout)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			return nil
		}
	}
}

var DefaultCloseTimeout = 60 * time.Second

type Option func(server *Server)

func NewServer(c *cron.Cron, opts ...Option) *Server {
	s := &Server{
		Cron:         c,
		closeTimeout: DefaultCloseTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

var _ reskit.Server = (*Server)(nil)
