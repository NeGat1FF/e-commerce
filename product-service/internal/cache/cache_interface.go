package cache

import (
	"context"
)

// Cache represents the cache client.
type Cache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string, res any) error
	Del(ctx context.Context, key string) error
}
