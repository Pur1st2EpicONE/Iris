package impl

import (
	"context"
	"database/sql"
	"errors"
)

// SaveVisit records a visit for the specified short URL and user agent.
// Errors are logged but do not propagate to the caller.
// If the short URL does not exist, no error is logged.
func (s *Service) SaveVisit(ctx context.Context, shortURL string, userAgent string) {
	err := s.storage.SaveVisit(ctx, shortURL, userAgent)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			s.logger.LogError("service — failed to save visit",
				err, "short link", shortURL, "user agent", userAgent, "layer", "service.impl")
		}
	}
}
