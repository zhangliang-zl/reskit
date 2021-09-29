package main

import (
	"context"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/component/redis"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/service"
	"github.com/zhangliang-zl/reskit/transport/httpx"
)

const (
	prjName = "test-prj"
)

func main() {

	logFac := logs.DefaultFactory
	rdsLogger, _ := logFac.Get("_redis")
	webLogger, _ := logFac.Get("_web")

	// init components
	kvstore, err := component.NewRedis(redis.Options{Addr: ":6379"}, rdsLogger)
	if err != nil {
		rdsLogger.Fatal(context.Background(), "error:%v", err)
	}

	// init service
	srv := service.NewHttpServer(httpx.Options{Addr: ":8096"}, webLogger)
	// httpx server setting
	engine := srv.Object().(*httpx.Engine)
	engine.UseMiddleware(httpx.PanicRecovery(webLogger))
	engine.AddRoute("GET", "/hello", func(ctx *httpx.Context) {
		ctx.Success(httpx.Map{"result": "world"})
	})

	app, _ := reskit.NewApp(
		reskit.WithName(prjName),
		reskit.WithEnv("dev"),

		// components
		reskit.WithComponent("redisIns1", kvstore),

		// services
		reskit.WithService("webServer1", srv),
	)
	appLogger, _ := logFac.Get("_app")
	err = app.Run()
	if err != nil {
		appLogger.Fatal(context.Background(), "app run :%v", err)
	}
}
