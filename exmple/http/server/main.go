package main

import (
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/component/redis"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/service"
	"github.com/zhangliang-zl/reskit/service/server/web"
	"log"
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
		log.Fatal(err)
	}

	// init service
	srv := service.NewHttpServer(web.Options{Addr: ":8096"}, webLogger)
	// web server setting
	engine := srv.Object().(*web.Engine)
	engine.UseMiddleware(web.PanicRecovery(webLogger))
	engine.AddRoute("GET", "/hello", func(ctx *web.Context) {
		ctx.Success(web.Map{"result": "world"})
	})

	app, _ := reskit.NewApp(
		reskit.WithName(prjName),
		reskit.WithEnv("dev"),

		// components
		reskit.WithComponent("redisIns1", kvstore),

		// services
		reskit.WithService("webServer1", srv),
	)

	log.Fatal(app.Run())
}
