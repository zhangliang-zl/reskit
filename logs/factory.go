package logs

import (
	"github.com/zhangliang-zl/reskit/logs/driver"
	"sync"
)

type Factory interface {
	Get(tag string) (Logger, error)
}

type factory struct {
	level   Level
	builder driver.WriterBuild
	sync.Map
}

func (f *factory) Get(tag string) (Logger, error) {
	v, ok := f.Load(tag)
	if !ok {
		record, err := f.builder(tag)
		if err != nil {
			return &logger{}, err
		}

		l := NewLogger(f.level, record)

		f.Store(tag, l)
		return l, err
	}

	return v.(Logger), nil
}

func NewFactory(level Level, writerBuild driver.WriterBuild) Factory {
	return &factory{
		level:   level,
		builder: writerBuild,
	}
}