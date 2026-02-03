package impl

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"errors"
)

func (s *Service) GetOriginalURL(ctx context.Context, link models.ShortLink) (string, error) {

	originalURL, err := s.cache.GetLink(ctx, link.ShortURL)
	if err == nil {
		s.logger.Debug("service — link fetched from cache", "short link", link.ShortURL, "layer", "service.impl")
		return originalURL, nil
	}

	originalURL, err = s.storage.GetOriginalURL(ctx, link.ShortURL)
	if err != nil {
		if errors.Is(err, errs.ErrLinkNotFound) {
			return "", err
		}
		s.logger.LogError("service — failed to get original url from DB", err, "short link", link.ShortURL, "layer", "service.impl")
		return "", err
	}

	if err := s.cache.SetLink(ctx, link.ShortURL, originalURL); err != nil {
		s.logger.LogError("service — failed to save link in cache", err, "short link", link.ShortURL, "layer", "service.impl")
	}

	s.logger.Debug("service — link fetched from DB", "short link", link.ShortURL, "layer", "service.impl")

	return originalURL, nil

}
