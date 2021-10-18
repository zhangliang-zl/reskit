package snowflake

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkIDWorkerGetId(b *testing.B) {
	// 表示最多16个节点，2的18次方 每秒最多生成 262144 ，70年内生成的最大ID不超过16位数字
	// 可满足大多数业务使用
	w, _ := NewWorker(0, 4, 18, time.Now().Unix()-86400*365*70)
	fmt.Println(w.NextID())
	for i := 0; i < b.N; i++ { // use b.N for looping
		w.NextID()
	}
}
