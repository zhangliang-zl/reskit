package beanstalked

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/mq"
	"time"
)

type Svc struct {
	queue    mq.Queue
	topic    string
	logger   logs.Logger
	stopChan chan bool
}

func (svc Svc) Stop() {
	svc.stopChan <- true
}

func (svc Svc) Serving(ctx context.Context, consumer mq.Consumer, fetchTimeout time.Duration) {

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

func NewBeanstalkSvc(topic string, queue mq.Queue, logger logs.Logger) mq.Svc {
	return Svc{
		topic:    topic,
		queue:    queue,
		logger:   logger,
		stopChan: make(chan bool),
	}
}
