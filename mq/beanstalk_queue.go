package mq

import (
	"context"
	"github.com/beanstalkd/go-beanstalk"
	"strconv"
	"sync"
	"time"
)

// Queue service based on beanstalk

type BeanstalkQueue struct {
	client *beanstalk.Conn
	tubeStore
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

func (q BeanstalkQueue) Push(ctx context.Context, topic string, body []byte, delay time.Duration) error {
	t := q.getTube(topic, q.client)
	_, err := t.Put(body, MinPriority, delay, MaxWorkingTTL)
	return err
}

func (q BeanstalkQueue) Fetch(ctx context.Context, topic string, timeout time.Duration) (body []byte, ask Ask, err error) {
	tube := beanstalk.NewTubeSet(q.client, topic)
	jobID, body, err := tube.Reserve(timeout)
	if err != nil {
		return
	}

	ask = q.buildAskFunc(ctx, jobID)
	return
}

// The ask function is used to process the execution result of the business
// The priority of handling failures is reduced until 2048
func (q BeanstalkQueue) buildAskFunc(ctx context.Context, jobID uint64) Ask {
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
		if pri <= MinPriority {
			pri = MinPriority
		}
		pri++
		if pri >= MaxPriority {
			pri = MaxPriority
		}

		return q.client.Release(jobID, pri, FailRetryDelay)
	}
}

func NewBeanstalkQueue(addr string) (Queue, error) {
	conn, err := beanstalk.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return BeanstalkQueue{
		client: conn,
	}, nil
}
