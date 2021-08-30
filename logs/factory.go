package logs

import (
	"github.com/zhangliang-zl/reskit/logs/driver"
	"sync"
)

type factory struct {
	level   LogLevel
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

func NewFactory(levelLabel string, writerBuild driver.WriterBuild) Factory {
	return &factory{
		level:   LevelVal(levelLabel),
		builder: writerBuild,
	}
}
