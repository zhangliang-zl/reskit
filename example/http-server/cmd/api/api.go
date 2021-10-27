package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/db"
	"github.com/zhangliang-zl/reskit/example/http-server/pkg/persist"
	"github.com/zhangliang-zl/reskit/example/http-server/pkg/route"
	"github.com/zhangliang-zl/reskit/redis"
	"github.com/zhangliang-zl/reskit/web"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
)

func main() {
	initLog()
	err := initPersist()
	if err != nil {
		persist.LogHelper.Fatal(err)
	}

	srv := web.New()
	srv.Middleware(web.SimpleTracing())
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

func initLog() {
	logger := log.With(log.DefaultLogger, "trace_id", TraceID())
	logger = log.With(logger, "span_id", SpanID())
	persist.LogHelper = log.NewHelper(log.With(logger, "project", "test-1"))
}

func TraceID() log.Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			return span.TraceID().String()
		}
		return fmt.Sprintf("%v", rand.Float32())
	}
}

func SpanID() log.Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
			return span.SpanID().String()
		}
		return fmt.Sprintf("%v", rand.Float32())
	}
}
