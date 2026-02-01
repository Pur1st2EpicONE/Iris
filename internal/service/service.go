package service

import (
	"Iris/internal/logger"
	"Iris/internal/models"
	"Iris/internal/repository"
	"Iris/internal/service/impl"
	"context"
)

type Service interface {
	ShortenLink(ctx context.Context, link models.Link) (string, error)
	GetOriginalURL(ctx context.Context, link models.ShortLink) (string, error)
	SaveVisit(ctx context.Context, shortURL string, userAgent string)
	GetAnalytics(ctx context.Context, shortURL string) (*models.VisitStats, error)
}

func NewService(logger logger.Logger, storage repository.Storage) Service {
	return impl.NewService(logger, storage)
}
