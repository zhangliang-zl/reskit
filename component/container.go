package component

import (
	"container/list"
	"context"
	"errors"
	"github.com/zhangliang-zl/reskit/logs"
	"sync"
)

type Container interface {
	Set(key string, instance Interface) error
	Get(key string) (interface{}, bool)
	Run() error
	Close() error
}

func NewContainer(logger logs.Logger, ctx context.Context) Container {
	return &container{
		components: make(map[string]Interface),
		enableSet:  true,
		logger:     logger,
		ctx:        ctx,
		entrySeq:   list.New(),
	}
}

type compWithKey struct {
	k string
	v Interface
}

type container struct {
	components map[string]Interface
	entrySeq   *list.List
	enableSet  bool
	logger     logs.Logger
	ctx        context.Context
	rw         sync.RWMutex
}

func (c *container) Set(k string, v Interface) error {
	c.rw.Lock()
	defer c.rw.Unlock()

	if !c.enableSet {
		return errors.New("Forbidden to call the Load() method after use RunFunc() ")
	}

	if _, ok := c.components[k]; ok {
		return errors.New("This Component Registered ")
	}

	c.components[k] = v
	c.entrySeq.PushBack(compWithKey{k: k, v: v})

	return nil
}

func (c *container) Get(k string) (interface{}, bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	comp, exist := c.components[k]
	if !exist {
		return nil, false
	}

	return comp.Instance(), true
}

func (c *container) Run() (err error) {
	c.enableSet = false

	for e := c.entrySeq.Front(); e != nil; e = e.Next() {
		comp := e.Value.(compWithKey)
		c.logger.Info(c.ctx, "component %s run", comp.k)
		if err = comp.v.Run(); err != nil {
			c.logger.Error(c.ctx, "component %s run err:  ->%s", comp.k, err.Error())
			return err
		}
	}

	return nil
}

func (c *container) Close() error {
	if len(c.components) > 0 {
		for e := c.entrySeq.Back(); e != nil; e = e.Prev() {
			comp := e.Value.(compWithKey)
			if err := comp.v.Close(); err != nil {
				c.logger.Error(c.ctx, "component %s close err %v", comp.k, err)
			} else {
				c.logger.Info(c.ctx, "component %s close", comp.k)
			}
		}
	}

	return nil
}
