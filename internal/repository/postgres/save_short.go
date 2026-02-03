package postgres

import (
	"context"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) SaveShort(ctx context.Context, id int64, shortLink string) error {

	query := `

	UPDATE links
	SET short_link = $1
	WHERE id = $2`

	_, err := s.db.ExecWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff,
	}, query, shortLink, id)

	return err

}
