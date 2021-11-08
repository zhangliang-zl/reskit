package web

import (
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/web/httperror"
	"runtime"
)

func Recovery(logger logs.Logger) HandlerFunc {
	return func(ctx *Context) {
		ctx.Set(logs.TraceIDKey, logs.NextTraceID())

		defer func() {
			// happen panic
			if r := recover(); r != nil {
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				logger.Error(ctx, "[panic] %v, stack: %s", r, string(buf[:n]))
				ctx.SendError(httperror.NewInternalError())
			}

			// error log
			lastErr := ctx.Errors.Last()
			if lastErr != nil {
				logger.Error(ctx, "[error]  %v ", lastErr)
			}
		}()

		ctx.Next()
	}
}
