package logs

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs/driver/stdout"
	"testing"
)

func TestLogger(t *testing.T) {
	f := NewFactory(LevelInfo, stdout.Driver())
	l, err := f.Get("haha")
	if err != nil {
		t.Fatal(err)
	}

	ctx := WithTraceID(context.Background())
	l.Error(ctx, "error message")
	l.Warn(ctx, "warn message")
}
