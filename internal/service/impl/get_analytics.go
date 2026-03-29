package impl

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"database/sql"
	"errors"
)

func (s *Service) GetAnalytics(ctx context.Context, groupBy string, shortURL string) (*models.VisitStats, error) {

	analytics, err := s.storage.GetAnalytics(ctx, groupBy, shortURL)
	if err == nil {
		return analytics, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrLinkNotFound
	}

	s.logger.LogError("service — failed to get analytics", err, "short link", shortURL, "layer", "service.impl")

	return nil, err

}
