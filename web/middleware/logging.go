package middleware

import (
	"bytes"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit/web"
	"io/ioutil"
	"strings"
)

func Logging(logger *log.Helper) web.HandlerFunc {
	return func(c *web.Context) {
		params := ""
		if c.ContentType() == "application/json" {
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			params += string(body)
		} else {
			c.Request.ParseForm()
			params += c.Request.PostForm.Encode()
		}

		msg := strings.ReplaceAll(params, "\n", "")
		if msg != "" {
			logger.Infof("[params] " + c.Request.RequestURI + " " + msg)
		}

		c.Next()
	}
}
