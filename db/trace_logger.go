package db

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

type traceLogger struct {
	logger *log.Helper
	gormlogger.Writer
	gormlogger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func (l *traceLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l traceLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger.Infof(msg, append([]interface{}{simplify(utils.FileWithLineNum())}, data...)...)
	}
}

func (l traceLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger.Warnf(msg, append([]interface{}{simplify(utils.FileWithLineNum())}, data...)...)
	}
}

func (l traceLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger.Errorf(l.errStr+msg, append([]interface{}{simplify(utils.FileWithLineNum())}, data...)...)
	}
}

func (l traceLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > gormlogger.Silent {
		elapsed := time.Since(begin)
		elapsedFloat := float64(elapsed.Nanoseconds()) / 1e6
		switch {
		case err != nil && l.LogLevel >= gormlogger.Error:
			errMsg := l.distinctError(err.Error())
			sql, rows := fc()
			if rows == -1 {
				l.logger.Errorf(l.traceErrStr, simplify(utils.FileWithLineNum()), errMsg, sql, "-", elapsedFloat)
			} else {
				l.logger.Errorf(l.traceErrStr, simplify(utils.FileWithLineNum()), errMsg, sql, rows, elapsedFloat)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
			sql, rows := fc()
			if rows == -1 {
				l.logger.Warnf(l.traceWarnStr, simplify(utils.FileWithLineNum()), sql, "-", elapsedFloat)
			} else {
				l.logger.Warnf(l.traceWarnStr, simplify(utils.FileWithLineNum()), sql, rows, elapsedFloat)
			}

		case l.LogLevel == gormlogger.Info:
			sql, rows := fc()
			if rows == -1 {
				l.logger.Infof(l.traceStr, simplify(utils.FileWithLineNum()), sql, "-", elapsedFloat)
			} else {
				l.logger.Infof(l.traceStr, simplify(utils.FileWithLineNum()), sql, rows, elapsedFloat)
			}
		}
	}
}

// Resolve gorm error logger repeat
func (l traceLogger) distinctError(errMsg string) string {
	if strings.LastIndex(errMsg, "; ErrorMsg") > 0 {
		repeatErrors := strings.Split(errMsg, "; ErrorMsg")
		return repeatErrors[0]
	}

	return errMsg
}

func simplify(file string) string {
	layers := strings.Split(file, "/")

	if len(layers) > 3 {
		return strings.Join(layers[len(layers)-3:], "/")
	}

	return file
}

func newTraceLogger(logger *log.Helper, config gormlogger.Config) gormlogger.Interface {
	var (
		infoStr      = "%s"
		warnStr      = "%s"
		errStr       = "%s"
		traceStr     = "%s %s  [rows:%v] [%.3fms]"
		traceWarnStr = "%s %s [rows:%v] [%.3fms]  [SLOW] "
		traceErrStr  = "%s %s %s [rows:%v] [%.3fms]  "
	)

	return &traceLogger{
		logger:       logger,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}
