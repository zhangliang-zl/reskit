package stdout

import (
	"github.com/zhangliang-zl/reskit/logs/driver"
	"io"
	"log"
)

func Driver() driver.WriterBuild {
	return func(tag string) (io.Writer, error) {
		return log.Writer(), nil
	}
}
