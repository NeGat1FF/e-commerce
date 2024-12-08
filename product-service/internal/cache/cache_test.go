package cache_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NeGat1FF/e-commerce/product-service/internal/cache"
	"github.com/NeGat1FF/e-commerce/product-service/internal/models"

	rd "github.com/testcontainers/testcontainers-go/modules/redis"
)

// func setupTestCache(t *testing.T) *redis.Client {
// 	opts, err := redis.ParseURL("redis://default:password@0.0.0.0:6379")
// 	require.NoError(t, err)

// 	return redis.NewClient(opts)
// }

var redisClient *redis.Client

func TestMain(m *testing.M) {
	redisContainer, err := rd.Run(context.Background(), "redis")
	if err != nil {
		panic(err)
	}

	str, err := redisContainer.Endpoint(context.Background(), "")
	if err != nil {
		panic(err)
	}

	fmt.Println(str)

	opts, err := redis.ParseURL(fmt.Sprintf("redis://:@%s", str))
	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(opts)

	code := m.Run()

	err = redisContainer.Terminate(context.Background())
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestRedisCache_Set(t *testing.T) {
	cache := cache.NewRedisClient(redisClient)

	pr := models.Product{
		ID:    1,
		Name:  "product",
		Price: 100,
	}

	err := cache.Set(context.Background(), "test:1", pr)
	require.NoError(t, err)

	var resProduct models.Product
	resString, err := redisClient.Get(context.Background(), "test:1").Result()
	require.NoError(t, err)

	err = json.Unmarshal([]byte(resString), &resProduct)
	require.NoError(t, err)

	assert.Equal(t, pr, resProduct)

	err = redisClient.Del(context.Background(), "test:1").Err()
	require.NoError(t, err)
}

func TestRedisCache_Get(t *testing.T) {
	cache := cache.NewRedisClient(redisClient)

	pr := models.Product{
		ID:    1,
		Name:  "product",
		Price: 100,
	}

	data, err := json.Marshal(pr)
	require.NoError(t, err)

	err = redisClient.Set(context.Background(), "test:1", data, 0).Err()
	require.NoError(t, err)

	var resProduct models.Product
	err = cache.Get(context.Background(), "test:1", &resProduct)
	require.NoError(t, err)

	assert.Equal(t, pr, resProduct)

	err = redisClient.Del(context.Background(), "test:1").Err()
	require.NoError(t, err)
}

func TestRedisCache_Del(t *testing.T) {
	cache := cache.NewRedisClient(redisClient)

	pr := models.Product{
		ID:    1,
		Name:  "product",
		Price: 100,
	}

	data, err := json.Marshal(pr)
	require.NoError(t, err)

	err = redisClient.Set(context.Background(), "test:1", data, 0).Err()
	require.NoError(t, err)

	err = cache.Del(context.Background(), "test:1")
	require.NoError(t, err)

	res, err := redisClient.Get(context.Background(), "test:1").Result()
	require.Error(t, err)
	assert.Equal(t, "", res)
}
