package mq

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs"
	"time"
)

type svcImpl struct {
	queue    Queue
	topic    string
	logger   logs.Logger
	stopChan chan bool
}

func (svc svcImpl) Stop() {
	svc.stopChan <- true
}

func (svc svcImpl) Serving(ctx context.Context, consumer Consumer, fetchTimeout time.Duration) {

loop:
	for {
		select {
		case <-svc.stopChan:
			return
			break
		default:
			data, callback, err := svc.queue.Fetch(ctx, svc.topic, fetchTimeout)

			if err != nil {
				if err.Error() == "reserve-with-timeout: timeout" {
					svc.logger.Info(ctx, "%s  no data .", svc.topic)
				} else {
					svc.logger.Error(ctx, " %s fetch error: %s", svc.topic, err.Error())
				}
				continue loop
			}

			err = consumer.Do(ctx, data)
			if err != nil {
				svc.logger.Error(ctx, "%s consumer do err %s", svc.topic, err.Error())
			}

			if err := callback(err); err != nil {
				svc.logger.Error(ctx, "%s callback error %s", svc.topic, err.Error())
			}
			break
		}
	}
}

func NewSvc(topic string, queue Queue, logger logs.Logger) Svc {
	return svcImpl{
		topic:    topic,
		queue:    queue,
		logger:   logger,
		stopChan: make(chan bool),
	}
}
