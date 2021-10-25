package syslog

import (
	"github.com/go-kratos/kratos/v2/log"
	"log/syslog"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, err := NewLogger(syslog.LOG_DEBUG|syslog.LOG_LOCAL6, "tag-1")
	if err != nil {
		panic(err)
	}
	helper:=log.NewHelper(logger)
	helper.Info("i am syslog logger testing ")
}
