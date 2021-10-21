package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit/web"
	"time"
)

func Speed(logger *log.Helper, slowThreshold int) web.HandlerFunc {
	return func(c *web.Context) {
		// register txID
		start := time.Now()
		defer func() {
			// slow query default 200 ms
			layout := "[speed] [%s] [%d] %s, usetime: %.3fms"
			elapsed := float64(time.Since(start).Microseconds()) / 1e3
			status := c.Writer.Status()
			params := []interface{}{c.Request.Method, status, c.Request.URL.Path, elapsed}

			if elapsed >= float64(slowThreshold) {
				logger.Warnf(layout+" [SLOW]", params...)
			} else {
				logger.Infof(layout, params...)
			}
		}()

		c.Next()
	}
}
