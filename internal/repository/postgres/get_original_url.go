package postgres

import (
	"Iris/internal/errs"
	"context"
	"database/sql"
	"errors"
)

func (s *Storage) GetOriginalURL(ctx context.Context, shortLink string) (string, error) {

	var originalURL string

	query := `

    SELECT original_link
    FROM links
    WHERE short_link = $1`

	err := s.db.QueryRowContext(ctx, query, shortLink).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrLinkNotFound
		}
		return "", err
	}

	return originalURL, nil

}
