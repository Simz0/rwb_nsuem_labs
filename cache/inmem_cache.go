package cache

import (
	"context"
	"sync"

	"labs/l0/database/repository"
)

type InMemCache struct {
	data sync.Map
}

func NewInMemCache() *InMemCache {
	return &InMemCache{}
}

func (c *InMemCache) Get(key string) (interface{}, bool) {
	return c.data.Load(key)
}

func (c *InMemCache) Set(key string, value interface{}) {
	c.data.Store(key, value)
}

func (c *InMemCache) Delete(key string) {
	c.data.Delete(key)
}

func (c *InMemCache) Preload(ctx context.Context, repo repository.OrderRepository) error {
	orders, err := repo.GetAllOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		c.Set(order.OrderUid, order)
	}

	return nil
}
