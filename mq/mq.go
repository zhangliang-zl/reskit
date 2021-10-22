package mq

import (
	"context"
	"time"
)

// Queue interface define  queue abstract
type Queue interface {
	Push(ctx context.Context, topic string, body []byte, delay time.Duration) error
	Fetch(ctx context.Context, topic string, timeout time.Duration) (body []byte, ask AskFunc, err error)
}

type AskFunc func(error) error

// Service interface define  queue abstract
type Service interface {
	Serving(ctx context.Context, topic string, queue Queue, consumer Consumer, fetchTimeout time.Duration)
	Stop() error
}

// Consumer  user self define
type Consumer interface {
	Do(ctx context.Context, raw []byte) error
}
