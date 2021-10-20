package db

import (
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Options struct {
	Dsn             string `json:"dsn"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int64  `json:"conn_max_lifetime"`
	ConnMaxIdleTime int64  `json:"conn_max_idle_time"`
	SlowThreshold   int64  `json:"slow_threshold"`
}

func New(opts Options, logger *log.Helper) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(opts.Dsn), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(opts.MaxOpenConns)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(opts.ConnMaxIdleTime) * time.Second)

	// default 100ms
	if opts.SlowThreshold == 0 {
		opts.SlowThreshold = 100
	}

	db.Logger = newTraceLogger(logger, gormlogger.Config{
		SlowThreshold: time.Duration(opts.SlowThreshold) * time.Millisecond,
		LogLevel:      gormlogger.Info,
	})
	return db, nil
}
