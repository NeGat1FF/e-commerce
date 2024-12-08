package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// Redis represents the Redis client.
type RedisCache struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client.
func NewRedisClient(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, 0).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, res any) error {
	resString, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(resString), &res)
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
