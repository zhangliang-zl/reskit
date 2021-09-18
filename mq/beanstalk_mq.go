package mq

import (
	"context"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/zhangliang-zl/reskit/logs"
	"strconv"
	"sync"
	"time"
)

// beanstalkQueue based on beanstalk
type beanstalkQueue struct {
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

func (q *beanstalkQueue) Push(ctx context.Context, topic string, body []byte, delay time.Duration) error {
	t := q.getTube(topic, q.client)
	_, err := t.Put(body, MinPriority, delay, MaxWorkingTTL)
	return err
}

func (q *beanstalkQueue) Fetch(ctx context.Context, topic string, timeout time.Duration) (body []byte, ask AskFunc, err error) {
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
func (q *beanstalkQueue) buildAskFunc(jobID uint64) AskFunc {
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

	return &beanstalkQueue{
		client: conn,
	}, nil
}

const (
	MinPriority    uint32 = 1024
	MaxPriority    uint32 = 2048
	MaxWorkingTTL         = time.Second * 120
	FailRetryDelay        = time.Second * 10
)

type beanstalkService struct {
	logger   logs.Logger
	stopChan chan bool
}

func (svc *beanstalkService) Stop() {
	svc.stopChan <- true
}

func (svc *beanstalkService) Serving(ctx context.Context, topic string, queue Queue, consumer Consumer, fetchTimeout time.Duration) {

loop:
	for {
		select {
		case <-svc.stopChan:
			return
		default:
			data, callback, err := queue.Fetch(ctx, topic, fetchTimeout)

			if err != nil {
				if err.Error() == "reserve-with-timeout: timeout" {
					svc.logger.Info(ctx, "%s  no data .", topic)
				} else {
					svc.logger.Error(ctx, " %s fetch error: %s", topic, err.Error())
				}
				continue loop
			}

			err = consumer.Do(ctx, data)
			if err != nil {
				svc.logger.Error(ctx, "%s consumer do err %s", topic, err.Error())
			}

			if err := callback(err); err != nil {
				svc.logger.Error(ctx, "%s callback error %s", topic, err.Error())
			}
			break
		}
	}
}

func NewBeanstalkService(logger logs.Logger) Service {
	return &beanstalkService{
		logger:   logger,
		stopChan: make(chan bool),
	}
}
