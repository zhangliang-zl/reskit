package dlock

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zhangliang-zl/reskit/redis"
	"sync"
	"testing"
	"time"
)

var factory Factory

func init() {

	var redisLogger = log.NewHelper(log.DefaultLogger, log.WithMessageKey("redis"))
	kvStore, _ := redis.New(redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}, redisLogger)

	var mutexLogger = log.NewHelper(log.DefaultLogger, log.WithMessageKey("mutex"))
	factory = NewRedisMutexFactory(mutexLogger, kvStore, "")
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
		t.Error("Concurrent dlock fail")
	}
}
