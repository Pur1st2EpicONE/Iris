package postgres

import (
	"context"

	"github.com/wb-go/wbf/retry"
)

// GetOriginalURL retrieves the original URL corresponding to a given short link.
// Returns the original URL as a string or an error if the short link does not exist.
func (s *Storage) GetOriginalURL(ctx context.Context, shortLink string) (string, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `

	SELECT original_link
	FROM links
	WHERE short_link = $1`,

		shortLink)
	if err != nil {
		return "", err
	}

	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		return "", err
	}

	return originalURL, nil

}
