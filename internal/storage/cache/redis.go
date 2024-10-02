package cache

import (
	"context"
	"fmt"

	"github.com/f24-cse535/apaxos/internal/config/storage"

	"github.com/redis/go-redis/v9"
)

// Cache is a module that use go-redis library to handle Redis operations.
type Cache struct {
	conn *redis.Client
}

// New creates a new cache instance with the given Redis configs.
// If the connection fails, it returns an error.
func New(cfg storage.Redis) (*Cache, error) {
	// open new redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	// check connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to open Redis connection: %v", err)
	}

	return &Cache{conn: rdb}, nil
}
