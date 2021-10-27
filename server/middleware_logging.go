package server

import (
	"bytes"
	"github.com/zhangliang-zl/reskit/logs"
	"io/ioutil"
	"strings"
)

func Logging(logger logs.Logger) HandlerFunc {
	return func(ctx *Context) {
		params := ""
		if ctx.ContentType() == "app/json" {
			body, _ := ioutil.ReadAll(ctx.Request.Body)
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			params += string(body)
		} else {
			_ = ctx.Request.ParseForm()
			params += ctx.Request.PostForm.Encode()
		}

		msg := strings.ReplaceAll(params, "\n", "")
		if msg != "" {
			logger.Info(ctx, "[params] uri:"+ctx.Request.RequestURI+",body:"+msg)
		}

		ctx.Next()
	}
}
