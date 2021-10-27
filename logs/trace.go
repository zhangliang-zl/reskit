package logs

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const traceIDKey = "trace_id"

func findTraceID(ctx context.Context) string {
	traceID := ctx.Value(traceIDKey)
	if traceID == nil {
		return ""
	}

	return traceID.(string)
}

func NextTraceID() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int31n(9999))
}
