// Package service defines the business logic layer for the Iris URL shortener.
// It provides an interface for operations such as shortening links, retrieving
// original URLs, recording visits, and fetching analytics.
package service

import (
	"Iris/internal/cache"
	"Iris/internal/logger"
	"Iris/internal/models"
	"Iris/internal/repository"
	"Iris/internal/service/impl"
	"context"
)

// Service defines the interface for the URL shortener service layer.
// It abstracts the core business logic from storage and caching details.
type Service interface {
	ShortenLink(ctx context.Context, link models.Link) (string, error)                             // ShortenLink creates a short URL for the given link.
	GetOriginalURL(ctx context.Context, link models.ShortLink) (string, error)                     // GetOriginalURL retrieves the original URL for a given short link.
	SaveVisit(ctx context.Context, shortURL string, userAgent string)                              // SaveVisit records a visit for a given short URL and user agent.
	GetAnalytics(ctx context.Context, groupBy string, shortURL string) (*models.VisitStats, error) // GetAnalytics retrieves visit statistics for a short URL.
}

// NewService creates a new Service implementation backed by the provided logger, cache, and storage.
func NewService(logger logger.Logger, cache cache.Cache, storage repository.Storage) Service {
	return impl.NewService(logger, cache, storage)
}
