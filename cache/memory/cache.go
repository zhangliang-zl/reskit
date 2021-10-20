package memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zhangliang-zl/reskit/cache"
	"sync"
	"time"
)

type Cache struct {
	size       int
	freeCache  *freecache.Cache
	mutexStore *sync.Map
}

func (m *Cache) getLocker(key string) *sync.Mutex {
	lock, ok := m.mutexStore.Load(key)
	if ok {
		return lock.(*sync.Mutex)
	}

	l := new(sync.Mutex)
	m.mutexStore.Store(key, l)
	return l
}

func (m *Cache) Get(ctx context.Context, key string, val interface{}) (bool, error) {
	b, err := m.freeCache.Get([]byte(key))
	if err != nil && err != freecache.ErrNotFound {
		return false, err
	}

	if err == freecache.ErrNotFound {
		return false, nil
	}

	err = msgpack.Unmarshal(b, val)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *Cache) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	data, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	return m.freeCache.Set([]byte(key), data, int(ttl.Seconds()))
}

func (m *Cache) Delete(ctx context.Context, key string) error {
	affected := m.freeCache.Del([]byte(key))
	if affected {
		return nil
	}

	return errors.New(fmt.Sprintf("Delete fail, No key %s", key))
}

func (m *Cache) GetOrSet(ctx context.Context, key string, val interface{}, ttl time.Duration, callback func() (interface{}, error)) (err error) {

	// When ttl<=0: callback() and return
	if ttl <= 0 {
		callbackVal, err := callback()
		if err == nil {
			cache.CopyObject(callbackVal, val)
		}
		return err
	}

	keyByte := []byte(key)

	// Get check Is Hit
	data, err := m.freeCache.Get(keyByte)

	if err == nil {
		return msgpack.Unmarshal(data, val)
	}

	if err != freecache.ErrNotFound {
		return err
	}

	// Key Not exist: callback() and set val
	lock := m.getLocker("Cache:" + key)

	lock.Lock()
	defer lock.Unlock()

	// Confirm again after getting the redislock
	data, err = m.freeCache.Get(keyByte)
	if len(data) > 0 {
		return msgpack.Unmarshal(data, val)
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

	err = m.freeCache.Set(keyByte, p, int(ttl.Seconds()))
	if err != nil {
		return err
	}

	return msgpack.Unmarshal(p, val)
}

func (c *Cache) init() {
	c.freeCache = freecache.NewCache(c.size)
}

type Option func(c *Cache)

// DefaultSize is 64M
const DefaultSize = 64 * 1024 * 1024

func Size(size int) Option {
	return func(c *Cache) {
		c.size = size
	}
}

func NewCache(opts ...Option) cache.Cache {
	c := &Cache{
		size:       DefaultSize,
		mutexStore: &sync.Map{},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.init()
	return c
}
