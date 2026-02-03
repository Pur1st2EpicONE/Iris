package postgres

import (
	"Iris/internal/errs"
	"context"
	"database/sql"
	"errors"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) SaveVisit(ctx context.Context, shortURL string, userAgent string) error {

	selectQuery := `
	
	SELECT id 
	FROM links 
	WHERE short_link = $1`

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff,
	}, selectQuery, shortURL)

	if err != nil {
		return err
	}

	var linkID int64
	if err := row.Scan(&linkID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrLinkNotFound
		}
		return err
	}

	insertQuery := `

	INSERT INTO visits (link_id, user_agent)
	VALUES ($1, $2)`

	_, err = s.db.ExecWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff,
	}, insertQuery, linkID, userAgent)

	return err

}
