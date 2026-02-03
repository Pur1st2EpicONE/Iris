package impl

import (
	"Iris/internal/cache"
	"Iris/internal/logger"
	"Iris/internal/repository"
)

type Service struct {
	logger  logger.Logger
	cache   cache.Cache
	storage repository.Storage
}

func NewService(logger logger.Logger, cache cache.Cache, storage repository.Storage) *Service {
	return &Service{logger: logger, cache: cache, storage: storage}
}
