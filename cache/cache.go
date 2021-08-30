package cache

import (
	"context"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, val interface{}) (bool, error)
	Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	GetOrSet(ctx context.Context, key string, val interface{}, ttl time.Duration, callback func() (interface{}, error)) (err error)
}

func copyObject(from interface{}, to interface{}) {
	b, _ := msgpack.Marshal(from)
	msgpack.Unmarshal(b, to)
}
