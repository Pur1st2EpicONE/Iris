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
}

func NewService(logger logger.Logger, storage repository.Storage) Service {
	return impl.NewService(logger, storage)
}
