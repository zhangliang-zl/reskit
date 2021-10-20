package beanstalkMQ

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)


type service struct {
	logger   *log.Helper
	stopChan chan bool
}

func (svc *service) Stop() {
	svc.stopChan <- true
}

func (svc *service) Serving(ctx context.Context, topic string, queue Queue, consumer Consumer, fetchTimeout time.Duration) {

loop:
	for {
		select {
		case <-svc.stopChan:
			return
		default:
			data, callback, err := queue.Fetch(ctx, topic, fetchTimeout)

			if err != nil {
				if err.Error() == "reserve-with-timeout: timeout" {
					svc.logger.Infof("%s  no data .", topic)
				} else {
					svc.logger.Errorf(" %s fetch error: %s", topic, err.Error())
				}
				continue loop
			}

			err = consumer.Do(ctx, data)
			if err != nil {
				svc.logger.Errorf("%s consumer do err %s", topic, err.Error())
			}

			if err := callback(err); err != nil {
				svc.logger.Errorf("%s callback error %s", topic, err.Error())
			}
			break
		}
	}
}

func NewBeanstalkService(logger *log.Helper) Service {
	return &service{
		logger:   logger,
		stopChan: make(chan bool),
	}
}
