package postgres

import (
	"context"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) SaveVisit(ctx context.Context, shortURL string, userAgent string) error {

	selectQuery := `
	
	SELECT id 
	FROM links 
	WHERE short_link = $1`

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy),
		selectQuery, shortURL)
	if err != nil {
		return err
	}

	var linkID int64
	if err := row.Scan(&linkID); err != nil {
		return err
	}

	insertQuery := `

	INSERT INTO visits (link_id, user_agent)
	VALUES ($1, $2)`

	_, err = s.db.ExecWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy),
		insertQuery, linkID, userAgent)

	return err

}
