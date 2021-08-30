package facade

import (
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/mq"
)

func RegisterBeanstalkQueue(addr string, id string) error {
	q, err := mq.NewBeanstalkQueue(addr)
	if err != nil {
		return err
	}

	instance := component.Make(q, nil, nil)
	return reskit.App().SetComponent(compQueue, id, instance)
}

func Queue(id string) mq.Queue {
	res, ok := reskit.App().Component(compQueue, id)

	if !ok {
		panic(compQueue + noRegister)
	}

	return res.(mq.Queue)
}
