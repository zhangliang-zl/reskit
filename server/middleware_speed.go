package server

import (
	"github.com/zhangliang-zl/reskit/logs"
	"time"
)

func Speed(logger logs.Logger, slowThreshold int) HandlerFunc {
	return func(ctx *Context) {
		// register txID
		start := time.Now()
		defer func() {
			// slow query default 200 ms
			layout := "[speed] [%s] [%d] %s, cost : %.3fms"
			elapsed := float64(time.Since(start).Microseconds()) / 1e3
			status := ctx.Writer.Status()
			params := []interface{}{ctx.Request.Method, status, ctx.Request.URL.Path, elapsed}

			if elapsed >= float64(slowThreshold) {
				logger.Warn(ctx, layout+" [SLOW]", params...)
			} else {
				logger.Info(ctx, layout, params...)
			}
		}()

		ctx.Next()
	}
}
