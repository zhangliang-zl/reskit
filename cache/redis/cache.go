package redis

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zhangliang-zl/reskit/cache"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/redislock"
	"time"
)

type Cache struct {
	client      *redis.Client
	keyPrefix   string
	logger      logs.Logger
	maxBuild    time.Duration // GetOrSet() The maximum time to build the cache, beyond which other threads can build
	lockFactory *redislock.Factory
}

func (c *Cache) Get(ctx context.Context, key string, val interface{}) (bool, error) {
	key = c.buildKey(key)
	return c.get(ctx, key, val)
}

func (c *Cache) get(ctx context.Context, key string, val interface{}) (exist bool, err error) {
	data, err := c.client.Get(ctx, key).Bytes()

	// Cache miss
	if err == redis.Nil {
		c.logger.Warn(ctx, "Cache Miss key %s", key)
		err = nil
		return
	}

	if err != nil {
		return
	}

	err = msgpack.Unmarshal(data, val)
	if err == nil {
		exist = true
	}
	return
}

func (c *Cache) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	key = c.buildKey(key)
	return c.set(ctx, key, val, ttl)
}

func (c *Cache) set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	b, err := msgpack.Marshal(val)
	err = c.client.Set(ctx, key, b, ttl).Err()
	if err != nil {
		c.logger.Error(ctx, "redis set %s error:%v", key, err)
	}
	return err
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	key = c.buildKey(key)
	err := c.delete(ctx, key)
	if err != nil {
		c.logger.Error(ctx, "Delete %s error:%v", key, err)
	}
	return err
}

func (c *Cache) GetOrSet(ctx context.Context, key string, valPointer interface{}, ttl time.Duration, callback func() (interface{}, error)) (err error) {
	err = c.getOrSet(ctx, key, valPointer, ttl, callback)
	if err != nil {
		c.logger.Error(ctx, "GetOrSet %s error:%v", key, err)
	}
	return err
}

func (c *Cache) getOrSet(ctx context.Context, key string, val interface{}, ttl time.Duration, callback func() (interface{}, error)) (err error) {
	// When ttl<=0: callback() and return
	if ttl <= 0 {
		callbackVal, err := callback()
		if err == nil {
			cache.CopyObject(callbackVal, val)
		}
		return err
	}

	// Key Exists:  return
	key = c.buildKey(key)
	exist, err := c.get(ctx, key, val)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	// Key Not exist: callback() and set val
	locker := c.lockFactory.New(key)

	_ = locker.Lock(ctx)
	defer locker.UnLock(ctx)

	// Confirm again after getting the redis
	exist, _ = c.get(ctx, key, val)
	if exist {
		return
	}

	callbackVal, err := callback()
	if err != nil {
		return err
	}

	// Set cache and return
	p, err := msgpack.Marshal(callbackVal)
	if err != nil {
		return err
	}

	err = c.client.Set(ctx, key, p, ttl).Err()
	if err != nil {
		c.logger.Error(ctx, "redis set %s error:%v", key, err)
	}

	return msgpack.Unmarshal(p, val)
}

func (c *Cache) delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) buildKey(key string) string {
	if len(key) < 256 {
		return c.keyPrefix + key
	}

	h := md5.New()
	h.Write([]byte(key))
	md5Key := hex.EncodeToString(h.Sum(nil))
	return c.keyPrefix + md5Key
}

func NewCache(client *redis.Client, lockFactory *redislock.Factory, opts ...Option) cache.Cache {
	c := &Cache{
		client:      client,
		logger:      DefaultLogger,
		keyPrefix:   DefaultKeyPrefix,
		lockFactory: lockFactory,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
