package impl

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"errors"
)

func (s *Service) GetAnalytics(ctx context.Context, shortURL string) (*models.VisitStats, error) {

	stats, err := s.storage.GetAnalytics(ctx, shortURL)
	if err == nil {
		return stats, nil
	}

	if errors.Is(err, errs.ErrLinkNotFound) {
		return nil, err
	}

	s.logger.LogError("service — failed to get analytics",
		err, "short link", shortURL, "layer", "service.impl")

	return nil, err

}
