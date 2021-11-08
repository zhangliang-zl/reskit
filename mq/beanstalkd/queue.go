package mq

import (
	"context"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/zhangliang-zl/reskit/mq"
	"strconv"
	"sync"
	"time"
)

// queue based on beanstalkd
type queue struct {
	client *beanstalk.Conn
	tubeStore
	opts *Options
}

type tubeStore struct {
	sync.Map
}

func (store tubeStore) getTube(topic string, client *beanstalk.Conn) *beanstalk.Tube {
	t, exist := store.Load(topic)
	if exist {
		return t.(*beanstalk.Tube)
	}

	tube := beanstalk.NewTube(client, topic)
	store.Store(topic, tube)
	return tube
}

func (q *queue) Push(_ context.Context, topic string, body []byte, delay time.Duration) error {
	t := q.getTube(topic, q.client)
	_, err := t.Put(body, q.opts.minPriority, delay, q.opts.maxWorkingTTL)
	return err
}

func (q *queue) Fetch(_ context.Context, topic string, timeout time.Duration) (body []byte, ask mq.AskFunc, err error) {
	tube := beanstalk.NewTubeSet(q.client, topic)
	jobID, body, err := tube.Reserve(timeout)
	if err != nil {
		return
	}

	ask = q.buildAskFunc(jobID)
	return
}

// The ask function is used to process the execution result of the business
// The priority of handling failures is reduced until 2048
func (q *queue) buildAskFunc(jobID uint64) mq.AskFunc {
	return func(err error) error {
		if err == nil {
			return q.client.Delete(jobID)
		}

		stat, err := q.client.StatsJob(jobID)
		if err != nil {
			return err
		}

		priInt, _ := strconv.Atoi(stat["pri"])
		pri := uint32(priInt)
		if pri <= q.opts.minPriority {
			pri = q.opts.minPriority
		}
		pri++
		if pri >= q.opts.minPriority {
			pri = q.opts.maxPriority
		}

		return q.client.Release(jobID, pri, q.opts.failRetryDelay)
	}
}

func NewQueue(addr string, opts ...Option) (mq.Queue, error) {
	conn, err := beanstalk.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	o := &Options{
		minPriority:    0,
		maxPriority:    1024,
		maxWorkingTTL:  time.Second * 120,
		failRetryDelay: time.Second * 10,
	}

	for _, opt := range opts {
		opt(o)
	}

	return &queue{
		client: conn,
		opts:   o,
	}, nil
}
