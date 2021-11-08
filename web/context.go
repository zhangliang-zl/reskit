package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangliang-zl/reskit/web/httperror"
	"net/http"
)

// Context Expand gin context functions
type Context struct {
	*gin.Context
}

func (c *Context) Success(obj interface{}) {
	c.Context.JSON(200, obj)
}

func (c *Context) SendError(err httperror.ErrorInterface) {
	sendError(c.Context, err)
}

func (c *Context) BadRequest(err error) {
	sendError(c.Context, httperror.NewBadRequest(err.Error()))
}

func (c *Context) InternalError(err error) {
	sendError(c.Context, httperror.NewInternalError(err.Error()))
}

func (c *Context) Unauthorized(err error) {
	sendError(c.Context, httperror.NewUnauthorized(err.Error()))
}

func (c *Context) Forbidden(err error) {
	sendError(c.Context, httperror.NewForbidden(err.Error()))
}

func sendError(c *gin.Context, err httperror.ErrorInterface) {
	code := err.GetCode()
	if code < 400 || code > 599 {
		code = http.StatusInternalServerError
	}

	_ = c.Error(err)
	c.Abort()
	c.JSON(err.GetCode(), map[string]interface{}{
		"error": map[string]interface{}{
			"code":    err.GetCode(),
			"message": err.Error(),
		},
	})
}
