// Package impl provides the concrete implementation of the Service interface.
package impl

import (
	"Iris/internal/cache"
	"Iris/internal/logger"
	"Iris/internal/repository"
)

// Service implements the business logic for URL shortening.
type Service struct {
	logger  logger.Logger      // logger is used for structured logging
	cache   cache.Cache        // cache stores recently accessed links
	storage repository.Storage // storage is the persistent backend
}

// NewService creates a new Service implementation.
func NewService(logger logger.Logger, cache cache.Cache, storage repository.Storage) *Service {
	return &Service{logger: logger, cache: cache, storage: storage}
}
