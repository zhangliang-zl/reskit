package web

import (
	"bytes"
	"github.com/go-kratos/kratos/v2/log"
	"io/ioutil"
	"strings"
)

func Logging(logger *log.Helper) HandlerFunc {
	return func(c *Context) {
		params := ""
		if c.ContentType() == "app/json" {
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			params += string(body)
		} else {
			_ = c.Request.ParseForm()
			params += c.Request.PostForm.Encode()
		}

		msg := strings.ReplaceAll(params, "\n", "")
		if msg != "" {
			logger.Infof("[params] " + c.Request.RequestURI + " " + msg)
		}

		c.Next()
	}
}
