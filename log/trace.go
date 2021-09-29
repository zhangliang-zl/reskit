package log

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const txID = "txID"

func GetTraceID(ctx context.Context) string {
	txID := ctx.Value(txID)

	if txID == nil {
		return ""
	}

	return txID.(string)
}

func WithTraceID(ctx context.Context) context.Context {
	id := ctx.Value(txID)

	if id == nil {
		id = NextTraceID()
		ctx = context.WithValue(ctx, txID, id)
	}
	return ctx
}

func NextTraceID() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int31n(9999))
}
