package logs

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const TraceID = "traceID"

func findTraceID(ctx context.Context) string {
	txID := ctx.Value(TraceID)

	if txID == nil {
		newID := NextTraceID()
		context.WithValue(ctx, TraceID, newID)
		return newID
	}

	return txID.(string)
}

func NextTraceID() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int31n(9999))
}
