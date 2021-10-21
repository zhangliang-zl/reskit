package redislock

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
	"time"
)

var ErrLockOvertime = errors.New("Locking failed , timeout ")

type Mutex struct {
	redisClient *redis.Client
	logger      *log.Helper

	key  string
	opts *Options

	// After this time, will no try lock(), default value 3*duration
	// lockOvertime = lockWaiting+ lockTime
	lockOvertime time.Time
	locked       bool
	lockID       string
}

func (m *Mutex) Lock(ctx context.Context) error {
	if m.locked {
		return nil
	}

	start := time.Now()
	m.lockOvertime = start.Add(m.opts.lockWaiting)
	m.lockID = uniqueID()

	for {
		cmd := m.redisClient.SetNX(ctx, m.key, m.lockID, m.opts.duration)
		if cmd.Err() == redis.ErrClosed {
			m.logger.Errorf("lock fail. %s, redis close", m.key)
			return redis.ErrClosed
		}

		if cmd.Err() != nil {
			m.logger.Errorf("lock fail. %s, redis error :%s", m.key, cmd.Err())
			continue
		}

		if cmd.Val() {
			m.logger.Infof("lock success. %s,  %s ,cost: %dms", m.key, m.lockID, time.Now().Sub(start).Milliseconds())
			m.locked = true

			return nil
		}

		time.Sleep(m.opts.retryInterval)

		if time.Now().UnixNano() > m.lockOvertime.UnixNano() {
			m.logger.Errorf("lock fail. %s err: lock timeout, cost: %dms", m.key, time.Now().Sub(start).Milliseconds())
			return ErrLockOvertime
		}
	}

}

func (m *Mutex) UnLock(ctx context.Context) {
	if !m.locked {
		return
	}

	script := "if redis.call('get', KEYS[1])==ARGV[1] then " +
		"	return redis.call('del', KEYS[1]) " +
		"else" +
		"	return 0 " +
		"end"
	cmd := m.redisClient.Eval(ctx, script, []string{m.key}, m.lockID)
	err := cmd.Err()
	if err != nil {
		m.logger.Errorf("unlock fail. %s, error :%v", m.key, err)
		return
	}

	m.locked = false

	m.logger.Infof("unlock %s, result: %v ", m.key, cmd.Val())
	return
}

func uniqueID() string {
	return strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
