package postgres

import (
	"context"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) SaveOriginal(ctx context.Context, originalURL string) (int64, error) {

	query := `

	INSERT INTO links (original_link)
	VALUES ($1)
	RETURNING id`

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy{Attempts: 2}, query, originalURL)
	if err != nil {
		return 0, err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
