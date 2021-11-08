package db

import (
	"github.com/zhangliang-zl/reskit/logs"
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
	Logger          logs.Logger
}

func New(dsn string, opts ...Option) (*gorm.DB, error) {

	o := &Options{
		maxOpenConns:    100,
		maxIdleConns:    100,
		connMaxIdleTime: 100 * time.Millisecond,
		connMaxLifetime: 300 * time.Second,
		slowThreshold:   300 * time.Second,
		disableAutoPing: true,
		Logger:          logs.DefaultLogger("_db"),
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
