package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhangliang-zl/reskit/server/httperror"
	"net/http"
	"sync"
)

type HandlerFunc func(*Context)

// Server Expand gin Server functions
type Server struct {
	*gin.Engine
	httpServer *http.Server
	opts       *Options
	mu         sync.Mutex
}

func (s *Server) AddRoute(method, relativePath string, handlers ...HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, 0)
	for _, h := range handlers {
		ginHandlers = append(ginHandlers, s.makeGinHandlerFunc(h))
	}
	s.Engine.Handle(method, relativePath, ginHandlers...)
}

func (s *Server) makeGinHandlerFunc(h func(*Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&Context{c})
	}
}

func (s *Server) UseMiddleware(h HandlerFunc) {
	s.Use(s.makeGinHandlerFunc(h))
}

func (s *Server) WrapPProf() {
	WrapPProf(s.Engine)
}

func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		s.opts.logger.Errorf("server shutdown error :%s", err.Error())
	}

	return err
}

func (s *Server) Run() error {
	s.mu.Lock()
	s.httpServer = &http.Server{
		Addr:         s.opts.address,
		Handler:      s.Engine,
		ReadTimeout:  s.opts.readTimeout,
		WriteTimeout: s.opts.writeTimeout,
	}
	s.mu.Unlock()
	s.opts.logger.Info("server start running")
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.opts.logger.Errorf("server run error :%s", err.Error())
	}
	return err
}

func New(opts ...Option) *Server {
	gin.SetMode(gin.ReleaseMode)

	o := &Options{
		address:      DefaultAddress,
		writeTimeout: DefaultWriteTimeout,
		readTimeout:  DefaultReadTimeout,
		logger:       DefaultLogger,
		middlewares:  DefaultMiddlewares,
	}

	for _, opt := range opts {
		opt(o)
	}

	s := &Server{
		Engine: gin.New(),
		opts:   o,
	}

	for _, middlewareFunc := range o.middlewares {
		s.UseMiddleware(middlewareFunc)
	}

	s.NoMethod(noMethod)
	s.NoMethod(noRoute)

	return s
}

func noRoute(c *gin.Context) {
	sendError(c, httperror.NewNotFound())
}

func noMethod(c *gin.Context) {
	sendError(c, httperror.NewMethodNotAllowed())
}
