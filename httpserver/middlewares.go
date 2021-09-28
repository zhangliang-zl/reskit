package httpserver

import (
	"bytes"
	"github.com/zhangliang-zl/reskit/logs"
	"io/ioutil"
	"runtime"
	"strings"
	"time"
)

func LogParams(logger logs.Logger) HandlerFunc {
	return func(c *Context) {
		params := ""
		if c.ContentType() == "app/json" {
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			params += string(body)
		} else {
			c.Request.ParseForm()
			params += c.Request.PostForm.Encode()
		}

		msg := strings.ReplaceAll(params, "\n", "")
		if msg != "" {
			logger.Info(c, "[params] "+c.Request.RequestURI+" "+msg)
		}

		c.Next()
	}
}

func LogSpeed(logger logs.Logger, slowThreshold int) HandlerFunc {
	return func(c *Context) {
		// register txID
		start := time.Now()
		defer func() {
			// slow query default 200 ms
			if slowThreshold == 0 {
				slowThreshold = 200
			}

			layout := "[speed] [%s] [%d] %s, usetime: %.3fms"
			elapsed := float64(time.Since(start).Microseconds()) / 1e3
			status := c.Writer.Status()
			params := []interface{}{c.Request.Method, status, c.Request.URL.Path, elapsed}

			if elapsed >= float64(slowThreshold) {
				logger.Warn(c, layout+" [SLOW]", params...)
			} else {
				logger.Info(c, layout, params...)
			}
		}()

		c.Next()
	}
}

func PanicRecovery(logger logs.Logger) HandlerFunc {
	return func(c *Context) {

		defer func() {
			// happen panic
			if r := recover(); r != nil {
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				logger.Error(c, "[panic] %v, stack: %s", r, string(buf[:n]))
				c.SendError(NewInternalError("Service internal error"))
			}

			// error log
			lastErr := c.Errors.Last()
			if lastErr != nil {
				logger.Error(c, "[error]  %v ", lastErr)
			}
		}()

		c.Next()
	}
}
