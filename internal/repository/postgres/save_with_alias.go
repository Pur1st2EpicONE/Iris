package postgres

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/wb-go/wbf/retry"
)

func (s *Storage) SaveWithAlias(ctx context.Context, link models.Link) error {

	query := `

	INSERT INTO links (original_link, short_link)
	VALUES ($1, $2)`

	_, err := s.db.ExecWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff,
	}, query, link.OriginalURL, link.Alias)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrAliasExists
		}
		return err
	}

	return nil

}
