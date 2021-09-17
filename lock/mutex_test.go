package lock

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/logs/driver/stdout"
	"github.com/zhangliang-zl/reskit/redis"
	"sync"
	"testing"
	"time"
)

var factory Factory

func init() {
	loggerFactory := logs.NewFactory(logs.LevelInfo, stdout.Driver())
	redisLogger, _ := loggerFactory.Get("redis")
	kvStore, _ := redis.New(redis.Options{
		Addr:     "localhost:6379",
		Password: "Password",
		DB:       0,
	}, redisLogger)

	mutexLogger, _ := loggerFactory.Get("mutex")
	factory = NewRedisMutexFactory(mutexLogger, kvStore)
}

func TestRedisMutex(t *testing.T) {
	concurrentNum := 10
	workTime := time.Millisecond * 150
	w := sync.WaitGroup{}
	w.Add(concurrentNum)

	start := time.Now()

	for i := 0; i < concurrentNum; i++ {
		go func() {
			ctx := context.Background()
			opts := Options{
				Key:         "testKey",
				Duration:    time.Millisecond * 100,
				LockWaiting: time.Second * 10,
				RenewTimes:  RenewForever,
			}
			locker := factory.New(opts)
			err := locker.Lock(ctx)
			defer locker.UnLock(ctx)
			if err != nil {
				t.Error(err)
			}
			time.Sleep(workTime)
			w.Done()
		}()
	}
	w.Wait()

	if time.Now().Sub(start) < time.Duration(int(workTime)*concurrentNum) {
		t.Error("Concurrent lock fail")
	}
}
