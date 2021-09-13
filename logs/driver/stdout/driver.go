package stdout

import (
	"fmt"
	"github.com/zhangliang-zl/reskit/logs/driver"
)

type writer struct{}

func (r writer) Record(m string) {
	fmt.Println(m)
}

func Driver() driver.WriterBuild {
	return func(tag string) (driver.Writer, error) {
		return writer{}, nil
	}
}
