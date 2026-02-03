package postgres

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"database/sql"
	"errors"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) GetAnalytics(ctx context.Context, shortURL string) (*models.VisitStats, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff}, `

		SELECT id
		FROM links
		WHERE short_link = $1`,

		shortURL)
	if err != nil {
		return nil, err
	}

	var linkID int64
	if err := row.Scan(&linkID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrLinkNotFound
		}
		return nil, err
	}

	stats := &models.VisitStats{
		ByUserAgent: make(map[string]int),
		ByDay:       make(map[string]int),
		ByMonth:     make(map[string]int),
	}

	rows, err := s.db.QueryWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff}, `
		
	SELECT
	COUNT(*) OVER() AS total_count, 
	to_char(visited_at, 'YYYY-MM-DD') AS day, 
	to_char(visited_at, 'YYYY-MM') AS month, 
	user_agent
	FROM visits
	WHERE link_id = $1`,

		linkID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {

		var total int
		var day, month, ua string

		if err := rows.Scan(&total, &day, &month, &ua); err != nil {
			return nil, err
		}

		stats.Count = total
		stats.ByDay[day]++
		stats.ByMonth[month]++
		stats.ByUserAgent[ua]++

	}

	return stats, nil

}
