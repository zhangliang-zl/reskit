package redis

import (
	"github.com/zhangliang-zl/reskit/cache/test"
	"github.com/zhangliang-zl/reskit/redis"
	"github.com/zhangliang-zl/reskit/redislock"
	"testing"
)

func TestCache(t *testing.T) {
	kvStore, _ := redis.New()
	lockFactory := redislock.NewFactory(kvStore)
	redisCache := NewCache(kvStore, lockFactory)
	test.AllCase(redisCache, t)
}
