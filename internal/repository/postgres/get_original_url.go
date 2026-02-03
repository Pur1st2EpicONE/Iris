package postgres

import (
	"Iris/internal/errs"
	"context"
	"database/sql"
	"errors"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) GetOriginalURL(ctx context.Context, shortLink string) (string, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff}, `

		SELECT original_link
		FROM links
		WHERE short_link = $1`,

		shortLink)
	if err != nil {
		return "", err
	}

	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrLinkNotFound
		}
		return "", err
	}

	return originalURL, nil

}
