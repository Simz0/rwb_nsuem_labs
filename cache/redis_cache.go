package cache

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"labs/l0/database/models"
	"labs/l0/database/repository"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // например, "localhost:6379"
	})

	return &RedisCache{client: rdb}
}

func (c *RedisCache) Get(key string) (interface{}, bool) {
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	var order models.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, false
	}

	return &order, true
}

func (c *RedisCache) Set(key string, value interface{}) {
	ctx := context.Background()
	order, ok := value.(*models.Order)
	if !ok {
		return
	}

	data, err := json.Marshal(order)
	if err != nil {
		return
	}

	c.client.Set(ctx, key, data, 0) // без TTL
}

func (c *RedisCache) Delete(key string) {
	ctx := context.Background()
	c.client.Del(ctx, key)
}

func (c *RedisCache) Preload(ctx context.Context, repo repository.OrderRepository) error {
	orders, err := repo.GetAllOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		c.Set(order.OrderUid, order)
	}

	return nil
}
