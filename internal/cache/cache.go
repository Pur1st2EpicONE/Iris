// Package cache provides abstractions for caching layer used in the application.
// It defines the Cache interface and exposes a constructor for concrete implementations.
package cache

import (
	"Iris/internal/cache/redis"
	"Iris/internal/config"
	"Iris/internal/logger"
	"context"
)

type Cache interface {
	SetLink(ctx context.Context, key string, value any) error
	GetLink(ctx context.Context, key string) (string, error)
	Close()
}

// Connect creates a new Cache instance (currently Redis) using the provided logger and configuration.
// It returns the initialized cache or an error if connection fails.
func Connect(logger logger.Logger, config config.Cache) (Cache, error) {
	return redis.Connect(logger, config)
}
