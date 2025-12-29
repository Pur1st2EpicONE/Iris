package service

import (
	"Iris/internal/logger"
	"Iris/internal/repository"
	"Iris/internal/service/impl"
)

type Service interface {
}

func NewService(logger logger.Logger, storage repository.Storage) Service {
	return impl.NewService(logger, storage)
}
