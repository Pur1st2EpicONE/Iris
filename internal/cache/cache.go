// Package cache provides abstractions for caching layer used in the application.
// It defines the Cache interface and exposes a constructor for concrete implementations.
package cache

import (
	"Iris/internal/cache/redis"
	"Iris/internal/config"
	"Iris/internal/logger"
	"context"
)

// Cache defines the interface for a caching layer used by the application.
type Cache interface {
	SetLink(ctx context.Context, key string, value any) error // SetLink sets the value for the given key in the cache.
	GetLink(ctx context.Context, key string) (string, error)  // GetLink retrieves the value for the given key from the cache.
	Close()                                                   // Close closes the cache connection and releases resources.
}

// Connect creates a new Cache instance (currently Redis) using the provided logger and configuration.
// It returns the initialized cache or an error if connection fails.
func Connect(logger logger.Logger, config config.Cache) (Cache, error) {
	return redis.Connect(logger, config)
}
