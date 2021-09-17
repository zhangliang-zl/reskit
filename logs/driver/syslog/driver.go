package syslog

import (
	"github.com/zhangliang-zl/reskit/logs/driver"
	"log/syslog"
)

type writer struct {
	w *syslog.Writer
}

func (r writer) Record(m string) {
	r.w.Info(m)
}

const defaultPriority = syslog.LOG_DEBUG | syslog.LOG_LOCAL6

func Driver(projectName string, priority syslog.Priority) driver.WriterBuild {
	return func(tag string) (driver.Writer, error) {
		tag = projectName + "/" + tag
		if priority == 0 {
			priority = defaultPriority
		}
		sysWriter, err := syslog.New(priority, tag)
		if err != nil {
			return writer{}, err
		}

		return writer{
			w: sysWriter,
		}, nil
	}
}
