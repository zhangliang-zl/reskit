package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhangliang-zl/reskit/logs"
	"net/http"
	"sync"
	"time"
)

// Engine : Expand gin Engine functions
type Engine struct {
	*gin.Engine
	httpServer   *http.Server
	addr         string
	readTimeout  int
	writeTimeout int
	logger       logs.Logger
	appCtx       context.Context
	mu           sync.Mutex
}

func (e *Engine) AddRoute(method, relativePath string, handlers ...HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, 0)
	for _, h := range handlers {
		ginHandlers = append(ginHandlers, e.makeGinHandlerFunc(h))
	}
	e.Engine.Handle(method, relativePath, ginHandlers...)
}

func (e *Engine) makeGinHandlerFunc(h func(*Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&Context{c})
	}
}

func (e *Engine) UseMiddleware(h HandlerFunc) {
	e.Use(e.makeGinHandlerFunc(h))
}

func (e *Engine) WrapPProf() {
	WrapPProf(e.Engine)
}

func (e *Engine) Stop(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.httpServer.Shutdown(context.Background())
}

func (e *Engine) Start(ctx context.Context) error {
	e.mu.Lock()
	e.httpServer = &http.Server{
		Addr:         e.addr,
		Handler:      e.Engine,
		ReadTimeout:  time.Duration(e.readTimeout) * time.Second,
		WriteTimeout: time.Duration(e.writeTimeout) * time.Second,
	}
	e.mu.Unlock()

	err := e.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		e.logger.Info(ctx, "web closed")
		return nil
	}

	return err
}

type HandlerFunc func(*Context)