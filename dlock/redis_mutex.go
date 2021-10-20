package dlock

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
	"time"
)

const (
	DefaultKeyPrefix     = "MutexLocker:"
	DefaultRetryInterval = 50 * time.Millisecond
	DefaultMaxRenewTimes = 30
	DefaultLocked        = 10 * time.Second
	DefaultLockWaiting   = 30 * time.Second

	RenewClose   = -2
	RenewForever = -1
)

type RedisMutex struct {
	redisClient *redis.Client
	logger      *log.Helper

	opts Options

	// After this time, will no try dlock(), default value 3*duration
	// lockOvertime = lockWaiting+ lockTime
	lockOvertime time.Time
	locked       bool
	lockID       string

	// Cancel renew channel
	cancel chan bool
}

func (m *RedisMutex) Lock(ctx context.Context) error {
	if m.locked {
		return nil
	}

	lockStart := time.Now()
	m.lockOvertime = lockStart.Add(m.opts.LockWaiting)
	m.lockID = uniqueID()

	for {
		cmd := m.redisClient.SetNX(ctx, m.opts.Key, m.lockID, m.opts.Duration)
		if cmd.Err() != nil {
			m.logger.Errorf("dlock fail. %s, redis error", m.opts.Key, cmd.Err())
		}

		if cmd.Val() {
			m.logger.Infof("dlock success. %s,  %s ,cost: %dms", m.opts.Key, m.lockID, time.Now().Sub(lockStart).Milliseconds())
			m.locked = true
			if m.opts.RetryInterval != RenewClose {
				go m.watchDog(ctx)
			}

			return nil
		}

		time.Sleep(m.opts.RetryInterval)

		if time.Now().UnixNano() > m.lockOvertime.UnixNano() {
			m.logger.Errorf("dlock fail. %s err: dlock timeout, cost: %dms", m.opts.Key, time.Now().Sub(lockStart).Milliseconds())
			return ErrLockOvertime
		}
	}

	m.locked = true

	return nil
}

//  Monitor whether the dlock expires and automatically renew
func (m *RedisMutex) watchDog(ctx context.Context) error {
	m.logger.Infof("dlock watch dog %s ", m.opts.Key)

	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("RedisMutex:watchDog panic %v", r)
			m.logger.Error(msg)
		}
	}()

	interval := m.opts.Duration / 3

	t := time.NewTicker(interval)

	script := "if redis.call('GET',KEYS[1]) == ARGV[1] then " +
		"	return redis.call('PEXPIRE',KEYS[1], ARGV[2]) " +
		"else " +
		"	return 0 " +
		"end"

	renew := 0

	for {
		select {
		case <-t.C:
			cmd := m.redisClient.Eval(ctx, script, []string{m.opts.Key}, m.lockID, strconv.Itoa(int(m.opts.Duration.Milliseconds())))
			if cmd.Err() != nil {
				m.logger.Errorf("lock renew fail. %s, error :%v", m.opts.Key, cmd.Err())
			} else {
				m.logger.Infof("lock renew success. %s, duration:%dms", m.opts.Key, int(m.opts.Duration.Milliseconds()))
			}

			renew++
			if m.opts.RenewTimes != RenewForever && renew > m.opts.RenewTimes {
				t.Stop()
				return nil
			}
		case <-m.cancel:
			m.logger.Infof("lock renew canceled . %s", m.opts.Key)
			t.Stop()
			return nil
		}
	}
}

func (m *RedisMutex) UnLock(ctx context.Context) {
	if !m.locked {
		return
	}

	if m.opts.RenewTimes != RenewClose {
		m.cancel <- true
	}

	script := "if redis.call('get', KEYS[1])==ARGV[1] then " +
		"	return redis.call('del', KEYS[1]) " +
		"else" +
		"	return 0 " +
		"end"
	cmd := m.redisClient.Eval(ctx, script, []string{m.opts.Key}, m.lockID)
	err := cmd.Err()
	if err != nil {
		m.logger.Errorf("unlock fail. %s, error :%v", m.opts.Key, err)
		return
	}

	m.locked = false

	m.logger.Infof("unlock %s, result: %v ", m.opts.Key, cmd.Val())
	return
}

type redisMutexFactory struct {
	redisClient *redis.Client
	logger      *log.Helper
	keyPrefix   string
}

func (factory redisMutexFactory) New(opts Options) Mutex {
	if opts.LockWaiting == 0 {
		opts.LockWaiting = DefaultLockWaiting
	}

	if opts.Duration == 0 {
		opts.Duration = DefaultLocked
	}

	if opts.RenewTimes == 0 {
		opts.RenewTimes = DefaultMaxRenewTimes
	}

	if opts.RenewTimes == 0 {
		opts.RetryInterval = DefaultRetryInterval
	}

	opts.Key = factory.keyPrefix + opts.Key

	return &RedisMutex{
		redisClient: factory.redisClient,
		logger:      factory.logger,
		cancel:      make(chan bool),
		opts:        opts,
	}
}

func NewRedisMutexFactory(logger *log.Helper, redisClient *redis.Client, keyPrefix string) Factory {
	if keyPrefix == "" {
		keyPrefix = DefaultKeyPrefix
	}

	return redisMutexFactory{
		logger:      logger,
		redisClient: redisClient,
		keyPrefix:   keyPrefix,
	}
}

func uniqueID() string {
	return strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
