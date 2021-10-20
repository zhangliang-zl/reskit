package beanstalkMQ

import (
	"context"
	"github.com/beanstalkd/go-beanstalk"
	"strconv"
	"sync"
	"time"
)

// queue based on beanstalk
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
	_, err := t.Put(body, q.opts.MinPriority, delay, q.opts.MaxWorkingTTL)
	return err
}

func (q *queue) Fetch(_ context.Context, topic string, timeout time.Duration) (body []byte, ask AskFunc, err error) {
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
func (q *queue) buildAskFunc(jobID uint64) AskFunc {
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
		if pri <= q.opts.MinPriority {
			pri = q.opts.MinPriority
		}
		pri++
		if pri >= q.opts.MinPriority {
			pri = q.opts.MaxPriority
		}

		return q.client.Release(jobID, pri, q.opts.FailRetryDelay)
	}
}

const (
	DefaultMinPriority    uint32 = 1024
	DefaultMaxPriority    uint32 = 2048
	DefaultMaxWorkingTTL         = time.Second * 120
	DefaultFailRetryDelay        = time.Second * 10
)

type Options struct {
	MinPriority    uint32
	MaxPriority    uint32
	MaxWorkingTTL  time.Duration
	FailRetryDelay time.Duration
}

type Option func(o *Options)

func NewQueue(addr string, opts ...Option) (Queue, error) {
	conn, err := beanstalk.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	o := &Options{
		MinPriority:    DefaultMinPriority,
		MaxPriority:    DefaultMaxPriority,
		MaxWorkingTTL:  DefaultMaxWorkingTTL,
		FailRetryDelay: DefaultFailRetryDelay,
	}

	return &queue{
		client: conn,
		opts:   o,
	}, nil
}
