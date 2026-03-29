package postgres

import (
	"context"

	"github.com/wb-go/wbf/retry"
)

// SaveShort updates the short link for a given link ID.
// Returns an error if the update fails.
func (s *Storage) SaveShort(ctx context.Context, id int64, shortLink string) error {

	query := `

	UPDATE links
	SET short_link = $1
	WHERE id = $2`

	_, err := s.db.ExecWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy),
		query, shortLink, id)

	return err

}
