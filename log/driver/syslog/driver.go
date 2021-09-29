package syslog

import (
	"github.com/zhangliang-zl/reskit/log/driver"
	"io"
	"log/syslog"
)

const defaultPriority = syslog.LOG_DEBUG | syslog.LOG_LOCAL6

func Driver(projectName string, priority syslog.Priority) driver.WriterBuild {
	return func(tag string) (io.Writer, error) {
		tag = projectName + "/" + tag
		if priority == 0 {
			priority = defaultPriority
		}
		sysWriter, err := syslog.New(priority, tag)
		return sysWriter, err
	}
}
