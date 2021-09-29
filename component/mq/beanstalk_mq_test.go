package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhangliang-zl/reskit/log"
	"github.com/zhangliang-zl/reskit/log/driver/stdout"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var q, _ = NewBeanstalkQueue("localhost:11300")
var loggerFactory = log.NewFactory(log.LevelWarn, stdout.Driver())
var logger, _ = loggerFactory.Get("mq")
var runTimes = 0

func TestQueueBasic(t *testing.T) {
	ctx := context.TODO()
	topic := "topic_basic"
	raw := Raw{
		Id:   1,
		Name: "basic",
		Data: "data_basic",
	}
	data, _ := json.Marshal(raw)

	err := q.Push(ctx, topic, data, time.Second*0)
	if err != nil {
		t.Error(err)
	}

	data, ask, err := q.Fetch(ctx, topic, time.Second*1)
	if err != nil {
		t.Error(err)
	}

	if err = ask(do(data)); err != nil {
		t.Error(err)
	}
}

func TestQueueSvc(t *testing.T) {
	ctx := context.TODO()
	topic := "topic_svc"
	wg := sync.WaitGroup{}
	c := &testConsumer{}

	qSvc := NewBeanstalkService(logger)
	wg.Add(1)
	go func() {
		qSvc.Serving(ctx, topic, q, c, time.Second*3)
		wg.Done()
	}()

	for i := 0; i < 50; i++ {
		in := Raw{
			Id:   i,
			Name: "beanstalkService",
			Data: "data_svc",
		}
		data, _ := json.Marshal(in)
		err := q.Push(ctx, topic, data, time.Second*0)
		if err != nil {
			t.Error(err)
		}
	}

	go func() {
		time.Sleep(time.Second * 4)
		qSvc.Stop()
	}()

	wg.Wait()
	fmt.Println(runTimes)
	if runTimes != 50 {
		t.Error("runTimes error")
	}
}

func do(data []byte) error {
	var raw Raw
	json.Unmarshal(data, &raw)

	// 50%  happen error
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Int63n(100)
	if n < 50 {
		return errors.New("happen error")
	}
	return nil
}

type Raw struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

type testConsumer struct {
	sync.Mutex
}

func (t *testConsumer) Do(ctx context.Context, v []byte) error {
	var raw Raw
	json.Unmarshal(v, &raw)
	t.Lock()
	runTimes += 1
	t.Unlock()
	return nil
}
