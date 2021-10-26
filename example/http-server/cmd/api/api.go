package main

import (
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/db"
	"github.com/zhangliang-zl/reskit/example/http-server/pkg/persist"
	"github.com/zhangliang-zl/reskit/example/http-server/pkg/route"
	"github.com/zhangliang-zl/reskit/redis"
	"github.com/zhangliang-zl/reskit/web"
)

func main() {
	err := initPersist()
	if err != nil {
		persist.LogHelper.Fatal(err)
	}

	srv := web.New()
	route.Init(srv)

	app := reskit.New(
		reskit.Servers(srv),
	)

	if err := app.Run(); err != nil {
		persist.LogHelper.Fatal(err)
	}
}

func initPersist() error {
	// init redis
	rds, err := redis.New()
	if err != nil {
		persist.LogHelper.Fatal(err)
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
