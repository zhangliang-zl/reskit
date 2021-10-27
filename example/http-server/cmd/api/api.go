package main

import (
	"context"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/db"
	"github.com/zhangliang-zl/reskit/example/http-server/pkg/persist"
	"github.com/zhangliang-zl/reskit/example/http-server/pkg/route"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/redis"
	"github.com/zhangliang-zl/reskit/web"
)

func main() {
	ctx := context.Background()
	err := initPersist()
	if err != nil {
		logs.DefaultLogger("app").Fatal(ctx, err.Error())
	}

	srv := web.New()
	srv.Middleware(web.SimpleTracing())
	route.Init(srv)

	app := reskit.New(
		reskit.Servers(srv),
	)

	if err := app.Run(); err != nil {
		logs.DefaultLogger("app").Fatal(ctx, err.Error())
	}
}

func initPersist() error {
	// init redis
	rds, err := redis.New()
	if err != nil {
		logs.DefaultLogger("app").Fatal(context.Background(), err.Error())
	}
	persist.KVStore = rds

	// init db
	t1DB, err := db.New("t1:111@tcp(127.0.0.1:3306)/t1?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		return err
	}

	persist.DB = t1DB

	return nil
}
