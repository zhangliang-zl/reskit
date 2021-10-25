package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit/server"
	"github.com/zhangliang-zl/reskit/server/httperror"
	"runtime"
)

func Recovery(logger *log.Helper) server.HandlerFunc {
	return func(c *server.Context) {
		defer func() {
			// happen panic
			if r := recover(); r != nil {
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				logger.Errorf("[panic] %v, stack: %s", r, string(buf[:n]))
				c.SendError(httperror.NewInternalError())
			}

			// error log
			lastErr := c.Errors.Last()
			if lastErr != nil {
				logger.Errorf("[error]  %v ", lastErr)
			}
		}()

		c.Next()
	}
}