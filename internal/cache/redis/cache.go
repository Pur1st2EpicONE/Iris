package redis

import (
	"Iris/internal/config"
	"Iris/internal/logger"
	"context"

	r "github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
)

// Cache implements the Cache interface using a Redis backend.
type Cache struct {
	client *r.Client     // underlying Redis client
	logger logger.Logger // structured logger
	config config.Cache  // cache configuration
}

// Connect establishes a connection to Redis using the provided logger and configuration.
// Returns a Cache instance or an error if the connection fails.
func Connect(logger logger.Logger, config config.Cache) (*Cache, error) {
	client, err := r.Connect(r.Options{
		Address:   config.Host + ":" + config.Port,
		Password:  config.Password,
		MaxMemory: config.MaxMemory,
		Policy:    config.Policy})
	if err != nil {
		return nil, err
	}
	return &Cache{client: client, logger: logger, config: config}, nil
}

// SetLink sets the value for a given key in Redis with expiration and retry strategy.
func (c *Cache) SetLink(ctx context.Context, key string, value any) error {
	return c.client.SetWithExpirationAndRetry(ctx, retry.Strategy{
		Attempts: c.config.RetryStrategy.Attempts,
		Delay:    c.config.RetryStrategy.Delay,
		Backoff:  c.config.RetryStrategy.Backoff},
		key, value, c.config.ExpirationTime)
}

// GetLink retrieves the value for a given key from Redis and refreshes its expiration.
func (c *Cache) GetLink(ctx context.Context, key string) (string, error) {
	if err := c.client.Expire(ctx, key, c.config.ExpirationTime); err != nil {
		return "", err
	}
	return c.client.Get(ctx, key)
}

// Close shuts down the Redis client and logs the outcome.
func (c *Cache) Close() {
	if err := c.client.Close(); err != nil {
		c.logger.LogError("redis — failed to close properly", err, "layer", "cache.redis")
	} else {
		c.logger.LogInfo("redis — cache closed", "layer", "cache.redis")
	}
}
