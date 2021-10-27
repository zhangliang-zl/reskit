package snowflake

import (
	"testing"
)

func BenchmarkIDWorkerGetId(b *testing.B) {
	w := NewWorker()
	for i := 0; i < b.N; i++ { // use b.N for looping
		w.NextID()
	}
}
