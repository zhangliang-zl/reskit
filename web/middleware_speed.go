package web

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

func Speed(logger *log.Helper, slowThreshold int) HandlerFunc {
	return func(c *Context) {
		// register txID
		start := time.Now()
		defer func() {
			// slow query default 200 ms
			layout := "[speed] [%s] [%d] %s, cost : %.3fms"
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
