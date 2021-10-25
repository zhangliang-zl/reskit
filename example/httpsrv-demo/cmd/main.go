package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit/application"
	"github.com/zhangliang-zl/reskit/server"
)

func main() {

	logHelper := log.NewHelper(log.With(log.DefaultLogger, "project", "test-1"))

	s := server.New()
	s.AddRoute("GET", "/get", func(c *server.Context) {
		c.JSON(200, map[string]string{
			"hello": "world",
		})
	})

	app := application.New(
		application.Servers(s),
	)

	logHelper.Fatal(app.Run())
}
