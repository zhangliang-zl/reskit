package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/web"
	"time"
)

func main() {

	logHelper := log.NewHelper(log.With(log.DefaultLogger, "project", "test-1"))

	srv := web.New(web.Address(":8081"))
	srv.AddRoute("GET", "/get", func(c *web.Context) {
		time.Sleep(3 * time.Second)
		c.JSON(200, map[string]string{
			"hello": "world",
		})
	})

	app := reskit.New(
		reskit.BeforeStart(
			func() error {
				fmt.Println("before start ")
				return nil
			}),
		reskit.AfterStart(
			func() error {
				fmt.Println("after start")
				return nil
			}),

		reskit.BeforeStop(
			func() error {
				fmt.Println("before stop")
				return nil
			}),

		reskit.AfterStop(
			func() error {
				fmt.Println("after stop")
				return nil
			}),

		reskit.Servers(srv),
	)

	if err := app.Run(); err != nil {
		logHelper.Fatal(err)
	}
}
