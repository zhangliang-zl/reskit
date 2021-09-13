package mq

import (
	"context"
	"time"
)

type Ask func(error) error

type Queue interface {
	Push(ctx context.Context, topic string, body []byte, delay time.Duration) error
	Fetch(ctx context.Context, topic string, timeout time.Duration) (body []byte, ask Ask, err error)
}

type Svc interface {
	Serving(ctx context.Context, consumer Consumer, fetchTimeout time.Duration)
	Stop()
}

type Consumer interface {
	Do(ctx context.Context, raw []byte) error
}

const (
	MinPriority    uint32 = 1024
	MaxPriority    uint32 = 2048
	MaxWorkingTTL         = time.Second * 120
	FailRetryDelay        = time.Second * 10
)
