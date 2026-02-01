package impl

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"errors"
)

func (s *Service) GetOriginalURL(ctx context.Context, link models.ShortLink) (string, error) {

	originalURL, err := s.storage.GetOriginalURL(ctx, link.ShortURL)
	if err == nil {
		return originalURL, nil
	}

	if errors.Is(err, errs.ErrLinkNotFound) {
		return "", err
	}

	s.logger.LogError("service — failed to get original url", err,
		"short link", link.ShortURL, "layer", "service.impl")

	return "", err

}
