package mq

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/mq"
	"time"
)

type service struct {
	logger   logs.Logger
	stopChan chan bool
}

func (svc *service) Stop() error {
	svc.stopChan <- true
	return nil
}

func (svc *service) Serving(ctx context.Context, topic string, queue mq.Queue, consumer mq.Consumer, fetchTimeout time.Duration) {

loop:
	for {
		select {
		case <-svc.stopChan:
			return
		default:
			data, callback, err := queue.Fetch(ctx, topic, fetchTimeout)

			if err != nil {
				if err.Error() == "reserve-with-timeout: timeout" {
					svc.logger.Info(ctx, "%s no data .", topic)
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

func NewService(logger logs.Logger) mq.Service {
	return &service{
		logger:   logger,
		stopChan: make(chan bool),
	}
}
