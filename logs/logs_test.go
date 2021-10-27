package logs

import (
	"context"
	"github.com/zhangliang-zl/reskit/logs/stdout"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(stdout.NewRecorder(), LevelInfo, "tag-1")
	logger.Info(context.Background(), "i am syslogs logger testing ")
}
