package lock

import (
	"context"
	"errors"
	"time"
)

type Mutex interface {
	Lock(ctx context.Context) error
	UnLock(ctx context.Context)
}

type Factory interface {
	New(Options) Mutex
}

type Options struct {
	Key           string
	Duration      time.Duration // default 10s
	LockWaiting   time.Duration // default 30s
	RetryInterval time.Duration // default 50ms
	RenewTimes    int           // default 30 times ; -1  forever ; -2 close
}


var ErrLockOvertime = errors.New("Locking failed , timeout ")
