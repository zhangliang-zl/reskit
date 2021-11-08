package transport

import (
	"context"
)

type Client interface {
	Get(ctx context.Context, uri string, param map[string]interface{}, result interface{}) ErrorInterface
	PostForm(ctx context.Context, uri string, data map[string]interface{}, result interface{}) ErrorInterface
	PostJson(ctx context.Context, uri string, jsonData interface{}, result interface{}) ErrorInterface
	Do(ctx context.Context, method string, uri string, result interface{}, params ...interface{}) ErrorInterface
}
