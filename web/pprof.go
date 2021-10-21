package web

import (
	"net/http/pprof"
	"strings"

	"github.com/gin-gonic/gin"
)

// WrapPProf Wrap adds several routes from package `net/http/pprof` to *gin.Engine object
func WrapPProf(router *gin.Engine) {
	wrapGroup(&router.RouterGroup)
}

// wrapGroup adds several routes from package `net/http/pprof` to *gin.RouterGroup object
func wrapGroup(router *gin.RouterGroup) {
	routers := []struct {
		Method  string
		Path    string
		Handler gin.HandlerFunc
	}{
		{"GET", "/debug/pprof/", indexHandler()},
		{"GET", "/debug/pprof/heap", heapHandler()},
		{"GET", "/debug/pprof/goroutine", goroutineHandler()},
		{"GET", "/debug/pprof/allocs", allocsHandler()},
		{"GET", "/debug/pprof/block", blockHandler()},
		{"GET", "/debug/pprof/threadcreate", threadCreateHandler()},
		{"GET", "/debug/pprof/cmdline", cmdlineHandler()},
		{"GET", "/debug/pprof/profile", profileHandler()},
		{"GET", "/debug/pprof/symbol", symbolHandler()},
		{"POST", "/debug/pprof/symbol", symbolHandler()},
		{"GET", "/debug/pprof/trace", traceHandler()},
		{"GET", "/debug/pprof/mutex", mutexHandler()},
	}

	basePath := strings.TrimSuffix(router.BasePath(), "/")
	var prefix string

	switch {
	case basePath == "":
		prefix = ""
	case strings.HasSuffix(basePath, "/debug"):
		prefix = "/debug"
	case strings.HasSuffix(basePath, "/debug/pprof"):
		prefix = "/debug/pprof"
	}

	for _, r := range routers {
		router.Handle(r.Method, strings.TrimPrefix(r.Path, prefix), r.Handler)
	}
}

// indexHandler will pass the call from /debug/pprof to pprof
func indexHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Index(ctx.Writer, ctx.Request)
	}
}

// heapHandler will pass the call from /debug/pprof/heap to pprof
func heapHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("heap").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// goroutineHandler will pass the call from /debug/pprof/goroutine to pprof
func goroutineHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("goroutine").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// allocsHandler will pass the call from /debug/pprof/allocs to pprof
func allocsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("allocs").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// blockHandler will pass the call from /debug/pprof/block to pprof
func blockHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("block").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// threadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof
func threadCreateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// cmdlineHandler will pass the call from /debug/pprof/cmdline to pprof
func cmdlineHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Cmdline(ctx.Writer, ctx.Request)
	}
}

// profileHandler will pass the call from /debug/pprof/profile to pprof
func profileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Profile(ctx.Writer, ctx.Request)
	}
}

// symbolHandler will pass the call from /debug/pprof/symbol to pprof
func symbolHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Symbol(ctx.Writer, ctx.Request)
	}
}

// traceHandler will pass the call from /debug/pprof/trace to pprof
func traceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Trace(ctx.Writer, ctx.Request)
	}
}

// mutexHandler will pass the call from /debug/pprof/mutex to pprof
func mutexHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("mutex").ServeHTTP(ctx.Writer, ctx.Request)
	}
}
