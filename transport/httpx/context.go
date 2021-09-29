package httpx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Context Expand gin context functions
type Context struct {
	*gin.Context
}

func (c *Context) Success(obj interface{}) {
	c.Context.JSON(200, obj)
}

func (c *Context) SendError(err ErrorInterface) {
	sendError(c.Context, err)
}

func (c *Context) BadRequest(err error) {
	sendError(c.Context, NewBadRequest(err.Error()))
}

func (c *Context) ClientError(httpCode int, err error) {
	sendError(c.Context, NewError(httpCode, err.Error()))
}

func (c *Context) InternalError(err error) {
	sendError(c.Context, NewInternalError(err.Error()))
}

func sendError(c *gin.Context, err ErrorInterface) {
	code := err.GetCode()
	if code < 400 || code > 599 {
		code = http.StatusInternalServerError
	}

	c.Error(err)
	c.Abort()
	c.JSON(err.GetCode(), map[string]interface{}{
		"error": map[string]interface{}{
			"code":    err.GetCode(),
			"message": err.Error(),
		},
	})
}
