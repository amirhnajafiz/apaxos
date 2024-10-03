package cache

import (
	"context"
	"fmt"
)

// Set a new key-value pair in Redis.
func (c Cache) Set(key string, value string) error {
	ctx := context.Background()
	key = fmt.Sprintf("%s-%s", c.prefix, key)

	response := c.conn.Set(ctx, key, value, 0)
	if err := response.Err(); err != nil {
		return fmt.Errorf("failed to set entry: %v", err)
	}

	return nil
}

// Get a value by its key in Redis.
func (c Cache) Get(key string) (string, error) {
	ctx := context.Background()
	key = fmt.Sprintf("%s-%s", c.prefix, key)

	response := c.conn.Get(ctx, key)
	if err := response.Err(); err != nil {
		return "", fmt.Errorf("failed to get entry: %v", err)
	}

	return response.String(), nil
}
