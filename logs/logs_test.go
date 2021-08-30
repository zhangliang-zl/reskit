package logs

import (
	"context"
	"fmt"
	"github.com/zhangliang-zl/reskit/logs/driver/stdout"
	"testing"
)

func TestLogger(t *testing.T) {
	f := NewFactory("warn", stdout.Driver())
	l, _ := f.Get("haha")
	fmt.Println(l.Level())
	l.Info(context.Background(), "haha")
}
