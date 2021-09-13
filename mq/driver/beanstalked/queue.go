package beanstalked

import (
	"context"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/zhangliang-zl/reskit/mq"
	"strconv"
	"sync"
	"time"
)

// queue srv based on beanstalk
type queue struct {
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

func (q queue) Push(ctx context.Context, topic string, body []byte, delay time.Duration) error {
	t := q.getTube(topic, q.client)
	_, err := t.Put(body, mq.MinPriority, delay, mq.MaxWorkingTTL)
	return err
}

func (q queue) Fetch(ctx context.Context, topic string, timeout time.Duration) (body []byte, ask mq.Ask, err error) {
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
func (q queue) buildAskFunc(ctx context.Context, jobID uint64) mq.Ask {
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
		if pri <= mq.MinPriority {
			pri = mq.MinPriority
		}
		pri++
		if pri >= mq.MaxPriority {
			pri = mq.MaxPriority
		}

		return q.client.Release(jobID, pri, mq.FailRetryDelay)
	}
}

func NewQueue(addr string) (queue, error) {
	conn, err := beanstalk.Dial("tcp", addr)

	if err != nil {
		return queue{}, err
	}

	return queue{
		client: conn,
	}, nil
}
