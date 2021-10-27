package server

import (
	"context"
)

func SimpleTracing() HandlerFunc {
	return func(c *Context) {
		context.WithValue(c, "trace_id", "123123")
		c.Next()
	}
}
