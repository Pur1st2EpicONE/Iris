package postgres

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"database/sql"
	"errors"
)

func (s *Storage) GetAnalytics(ctx context.Context, shortURL string) (*models.VisitStats, error) {

	var linkID int64
	err := s.db.QueryRowContext(ctx, `
		
	SELECT id
	FROM links
	WHERE short_link = $1`,

		shortURL).Scan(&linkID)
	if err != nil {
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

	rows, err := s.db.QueryContext(ctx, `
		
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
	defer rows.Close()

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
