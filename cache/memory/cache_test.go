package memory

import (
	"github.com/zhangliang-zl/reskit/cache/test"
	"testing"
)

func TestCache(t *testing.T) {
	test.AllCase(NewCache(), t)
}
