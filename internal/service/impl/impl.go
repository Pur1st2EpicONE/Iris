package impl

import (
	"Iris/internal/logger"
	"Iris/internal/models"
	"Iris/internal/repository"
	"context"
)

type Service struct {
	logger  logger.Logger
	storage repository.Storage
}

func NewService(logger logger.Logger, storage repository.Storage) *Service {
	return &Service{logger: logger, storage: storage}
}

func (s *Service) ShortenLink(ctx context.Context, link models.Link) (string, error) {

	if err := validateLink(link); err != nil {
		return "", err
	}

	

	return "", nil
}
