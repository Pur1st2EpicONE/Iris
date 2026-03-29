package postgres

import (
	"context"

	"github.com/wb-go/wbf/retry"
)

// SaveOriginal inserts a new original URL into the links table and returns its generated ID.
// Returns the ID of the newly created link or an error if the operation fails.
func (s *Storage) SaveOriginal(ctx context.Context, originalURL string) (int64, error) {

	query := `

	INSERT INTO links (original_link)
	VALUES ($1)
	RETURNING id`

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy),
		query, originalURL)
	if err != nil {
		return 0, err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil

}
