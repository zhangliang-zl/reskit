package server

import (
	"bytes"
	"github.com/zhangliang-zl/reskit/logs"
	"io/ioutil"
	"strings"
)

func RequestParams(logger logs.Logger) HandlerFunc {
	return func(ctx *Context) {
		params := ""
		if ctx.ContentType() == "application/json" {
			body, _ := ioutil.ReadAll(ctx.Request.Body)
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			params += string(body)
		} else {
			_ = ctx.Request.ParseForm()
			params += ctx.Request.PostForm.Encode()
		}

		body := strings.ReplaceAll(params, "\n", "")
		msg := "[request] uri:" + ctx.Request.RequestURI
		if body != "" {
			msg += ",body:" + body
		}

		logger.Info(ctx, msg)
		ctx.Next()
	}
}
