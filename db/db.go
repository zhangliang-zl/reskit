package db

import (
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Options struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	slowThreshold   time.Duration
	disableAutoPing bool
	Logger          *log.Helper
}

func New(dsn string, opts ...Option) (*gorm.DB, error) {

	o := &Options{
		maxOpenConns:    DefaultMaxOpenConns,
		maxIdleConns:    DefaultMaxIdleConns,
		connMaxIdleTime: DefaultConnMaxIdleTime,
		connMaxLifetime: DefaultConnMaxLifetime,
		slowThreshold:   DefaultSlowThreshold,
		disableAutoPing: DefaultDisableAutoPing,
		Logger:          DefaultLogger,
	}

	for _, opt := range opts {
		opt(o)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableAutomaticPing: o.disableAutoPing})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(o.maxOpenConns)
	sqlDB.SetMaxIdleConns(o.maxIdleConns)
	sqlDB.SetConnMaxLifetime(o.connMaxIdleTime)

	db.Logger = newTraceLogger(o.Logger, gormlogger.Config{
		SlowThreshold: o.slowThreshold,
		LogLevel:      gormlogger.Info,
	})

	return db, nil
}
