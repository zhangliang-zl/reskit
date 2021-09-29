package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangliang-zl/reskit/logs"
)

type Options struct {
	Addr          string `json:"addr"`
	ReadTimeout   int    `json:"read_timeout"`   // default 60
	WriteTimeout  int    `json:"write_timeout"`  // default 60
	SlowThreshold int    `json:"slow_threshold"` // Slow query default 200 ms
	LogLevel      string `json:"logLevel"`
}

func NewServer(opts Options, logger logs.Logger) *Engine {
	if opts.LogLevel != "" {
		logger.SetLevel(logs.LevelVal(opts.LogLevel))
	}

	gin.SetMode(gin.ReleaseMode)

	e := &Engine{
		Engine: gin.New(),
	}

	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = 60
	}

	if opts.WriteTimeout == 0 {
		opts.WriteTimeout = 60
	}
	e.addr = opts.Addr
	e.readTimeout = opts.ReadTimeout
	e.writeTimeout = opts.WriteTimeout
	e.logger = logger

	e.NoMethod(noMethod)
	e.NoMethod(noRoute)
	e.UseMiddleware(PanicRecovery(logger))
	e.UseMiddleware(LogSpeed(logger, opts.SlowThreshold))

	return e
}

func noRoute(c *gin.Context) {
	sendError(c, NewNotFound("404 not found"))
}

func noMethod(c *gin.Context) {
	sendError(c, NewMethodNotAllowed("405 method not allowed"))
}
