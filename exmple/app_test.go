package exmple

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/facade"
	"github.com/zhangliang-zl/reskit/logs"
	kvstore "github.com/zhangliang-zl/reskit/redis"
	"log"
	"testing"
	"time"
)

type testService struct {
}

func (t testService) ID() string {
	return "exmple-job"
}

func (t testService) Object() interface{} {
	return "exmple-job-object"
}

func (t testService) Start(_ context.Context) error {

	time.Sleep(1 * time.Second)
	redisClient := facade.Redis("main")
	redisClient.Set(context.Background(), "a", "1", time.Second*10)
	val := redisClient.Get(context.Background(), "a").String()
	fmt.Println(val)

	time.Sleep(3 * time.Second)
	return errors.New("i am error")
}

func (t testService) Stop(_ context.Context) error {
	return nil
}

func TestApp_Close(t *testing.T) {

	instance, _ := facade.NewApp(
		app.WithName("zhang"),
		app.WithVersion("0.0.0.1"),
		app.WithEnv("exmple"),
		app.WithSignal(nil),
		app.WithContext(nil),
		app.WithLoggerFactory(logs.DefaultFactory),
	)

	if err := facade.RegisterRedis(kvstore.Options{Addr: ":6379"}, "main"); err != nil {
		log.Fatal(err)
	}

	svc1 := testService{}
	instance.RegisterService(svc1, "svc1")
	instance.RegisterService(svc1, "svc2")

	err := instance.Run()
	fmt.Println(err)
}
