package snowflake

import (
	"sync"
	"time"
)

type Worker struct {
	id         int64
	workerBits int64
	numberBits int64 // 每秒最大位数
	epoch      int64 // system start use time

	numberMax   int64
	timeShift   int64
	workerShift int64

	mu        sync.Mutex
	timestamp int64
	seq       int64
}

func (w *Worker) NextID() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UnixNano() / 1e9
	if w.timestamp == now {
		w.seq++

		if w.seq > w.numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e9
			}
		}
	} else {
		w.seq = 0
		w.timestamp = now
	}

	return (now-w.epoch)<<w.timeShift | (w.id << w.workerShift) | (w.seq)
}

// NewWorker is Worker Constructor
func NewWorker(opts ...Option) *Worker {

	w := &Worker{
		id:         0,
		workerBits: 5,          // 32个节点
		numberBits: 17,         // 每秒最多生产2的17次方 (131072) 个ID
		epoch:      1621481706, // 2021-05-20 11:35:06
	}

	for _, opt := range opts {
		opt(w)
	}

	var (
		workerMax   int64 = -1 ^ (-1 << w.workerBits)
		numberMax   int64 = -1 ^ (-1 << w.numberBits)
		timeShift         = w.workerBits + w.numberBits
		workerShift       = w.numberBits
		id                = w.id % workerMax
	)

	w.id = id
	w.timeShift = timeShift
	w.workerShift = workerShift
	w.numberMax = numberMax
	return w
}

type Option func(worker *Worker)

func WorkerIDBits(workerBits int64) Option {
	return func(worker *Worker) {
		worker.workerBits = workerBits
	}
}

func Id(id int64) Option {
	return func(worker *Worker) {
		worker.id = id
	}
}

func NumberBits(numberBits int64) Option {
	return func(worker *Worker) {
		worker.numberBits = numberBits
	}
}

func Epoch(e int64) Option {
	return func(worker *Worker) {
		worker.epoch = e
	}
}
