package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 5  // 0-31
	numberBits  uint8 = 17 // 0-131071
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift         = workerBits + numberBits
	workerShift       = numberBits
	epoch       int64 = 1621481706
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func (w *Worker) NextID() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UnixNano() / 1e9
	if w.timestamp == now {
		w.number++

		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e9
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}

	return (now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number)
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("Worker ID excess of quantity ")
	}
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}
