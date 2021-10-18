package snowflake

import (
	"sync"
	"time"
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64

	numberBits  int64 // 每秒最大位数
	numberMax   int64
	timeShift   int64
	workerShift int64
	epoch       int64 // system start use time
}

func (w *Worker) NextID() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UnixNano() / 1e9
	if w.timestamp == now {
		w.number++

		if w.number > w.numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e9
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}

	return (now-w.epoch)<<w.timeShift | (w.workerId << w.workerShift) | (w.number)
}

// NewWorker is Worker Constructor
// workerBits workerID 位数
// numberBits 每秒可容纳的ID数量
// epoch  表示系统开始使用该worker的时间戳
func NewWorker(workerId, workerBits, numberBits, epoch int64) (*Worker, error) {
	var (
		workerMax   int64 = -1 ^ (-1 << workerBits)
		numberMax   int64 = -1 ^ (-1 << numberBits)
		timeShift         = workerBits + numberBits
		workerShift       = numberBits
	)

	workerId = workerId % workerMax

	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,

		numberMax:   numberMax,
		timeShift:   timeShift,
		workerShift: workerShift,
		epoch:       epoch,
	}, nil
}
