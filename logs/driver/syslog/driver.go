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

func Driver(projectName string) driver.WriterBuild {
	return func(tag string) (driver.Writer, error) {
		tag = projectName + "/" + tag
		sysWriter, err := syslog.New(syslog.LOG_DEBUG|syslog.LOG_LOCAL6, tag)
		if err != nil {
			return writer{}, err
		}

		return writer{
			w: sysWriter,
		}, nil
	}
}
