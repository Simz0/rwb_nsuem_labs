package cache

import (
	"context"

	"labs/l0/database/repository"
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Preload(ctx context.Context, repo repository.OrderRepository) error
}
