package facade

import (
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/mq"
	beanstalked2 "github.com/zhangliang-zl/reskit/mq/driver/beanstalked"
)

func RegisterBeanstalkQueue(addr string, id string) error {
	q, err := beanstalked2.NewQueue(addr)
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
