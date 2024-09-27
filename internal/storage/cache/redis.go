package cache

import (
	"context"
	"fmt"

	"github.com/f24-cse535/apaxos/internal/config/storage"

	"github.com/redis/go-redis/v9"
)

// Cache is a module that handles redis operations.
type Cache struct {
	conn *redis.Client
}

func NewCache(cfg storage.RedisConfig) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to open Redis connection: %v", err)
	}

	return &Cache{conn: rdb}, nil
}
