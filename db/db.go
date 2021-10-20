package db

import (
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Options struct {
	MaxOpenConns         int
	MaxIdleConns         int
	ConnMaxLifetime      time.Duration
	ConnMaxIdleTime      time.Duration
	SlowThreshold        time.Duration
	DisableAutomaticPing bool
	Logger               *log.Helper
}

var (
	DefaultMaxOpenConns         = 100
	DefaultMaxIdleConns         = 100
	DefaultSlowThreshold        = 100 * time.Millisecond
	DefaultConnMaxLifetime      = 300 * time.Second
	DefaultConnMaxIdleTime      = 300 * time.Second
	DefaultDisableAutomaticPing = true
	DefaultLogger               = log.NewHelper(log.DefaultLogger, log.WithMessageKey("db"))
)

type Option func(options *Options)

func New(dsn string, opts ...Option) (*gorm.DB, error) {

	o := &Options{
		MaxOpenConns:         DefaultMaxOpenConns,
		MaxIdleConns:         DefaultMaxIdleConns,
		ConnMaxIdleTime:      DefaultConnMaxIdleTime,
		ConnMaxLifetime:      DefaultConnMaxLifetime,
		SlowThreshold:        DefaultSlowThreshold,
		DisableAutomaticPing: DefaultDisableAutomaticPing,
		Logger:               DefaultLogger,
	}

	for _, opt := range opts {
		opt(o)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableAutomaticPing: o.DisableAutomaticPing})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(o.MaxOpenConns)
	sqlDB.SetMaxIdleConns(o.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(o.ConnMaxIdleTime)

	db.Logger = newTraceLogger(o.Logger, gormlogger.Config{
		SlowThreshold: o.SlowThreshold,
		LogLevel:      gormlogger.Info,
	})
	return db, nil
}
