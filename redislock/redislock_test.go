package redislock

import (
	"context"
	"github.com/zhangliang-zl/reskit/redis"
	"sync"
	"testing"
	"time"
)

var factory *Factory

func init() {
	kvStore, _ := redis.New()
	factory = NewFactory(kvStore)
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
			locker := factory.New("testKey",
				LockTime(time.Second),
				LockWaiting(time.Second*30),
			)
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
