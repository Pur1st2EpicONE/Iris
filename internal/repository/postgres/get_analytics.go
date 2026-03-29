package postgres

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// GetAnalytics retrieves visit statistics for a given short URL.
// The groupBy parameter supports "day", "month", "user_agent", or "" for total count.
func (s *Storage) GetAnalytics(ctx context.Context, groupBy string, shortURL string) (*models.VisitStats, error) {

	linkID, err := s.getLinkID(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	if groupBy == "" {
		return s.getTotalCount(ctx, linkID)
	}
	return s.getStats(ctx, linkID, groupBy)

}

// getLinkID retrieves the database ID for a given short URL.
func (s *Storage) getLinkID(ctx context.Context, shortURL string) (int64, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `
		
	SELECT id
	FROM links
	WHERE short_link = $1
	
	`, shortURL)
	if err != nil {
		return 0, err
	}

	var linkID int64
	if err := row.Scan(&linkID); err != nil {
		return 0, err
	}

	return linkID, nil

}

// getTotalCount retrieves all visit records for the given link ID without grouping.
func (s *Storage) getTotalCount(ctx context.Context, linkID int64) (*models.VisitStats, error) {

	rows, err := s.db.QueryWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `
		
	SELECT visited_at, user_agent
	FROM visits
	WHERE link_id = $1
	ORDER BY visited_at
	
	`, linkID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	stats := &models.VisitStats{Data: []models.VisitEntry{}}

	for rows.Next() {
		var visitedAt, userAgent string
		if err := rows.Scan(&visitedAt, &userAgent); err != nil {
			return nil, err
		}
		stats.Data = append(stats.Data, models.VisitEntry{
			Key:       "",
			UserAgent: userAgent,
			Time:      visitedAt,
			Count:     1,
		})
		stats.Count++
	}

	return stats, nil

}

// getStats retrieves visit records for the given link ID grouped by the specified parameter.
// Supported groupBy values are "day", "month", or "user_agent".
func (s *Storage) getStats(ctx context.Context, linkID int64, groupBy string) (*models.VisitStats, error) {

	var groupExpr string

	switch groupBy {
	case "day":
		groupExpr = "to_char(visited_at, 'YYYY-MM-DD')"
	case "month":
		groupExpr = "to_char(visited_at, 'YYYY-MM')"
	case "user_agent":
		groupExpr = "user_agent"
	default:
		return nil, errs.ErrInvalidGroupBy
	}

	query := fmt.Sprintf(`

	SELECT visited_at, user_agent, %s AS key
	FROM visits
	WHERE link_id = $1
	ORDER BY visited_at
	
	`, groupExpr)

	rows, err := s.db.QueryWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), query, linkID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	stats := &models.VisitStats{Data: []models.VisitEntry{}}

	for rows.Next() {
		var visitedAt string
		var userAgent string
		var key string
		if err := rows.Scan(&visitedAt, &userAgent, &key); err != nil {
			return nil, err
		}
		stats.Data = append(stats.Data, models.VisitEntry{
			Key:       key,
			UserAgent: userAgent,
			Time:      visitedAt,
			Count:     1,
		})
		stats.Count++
	}

	return stats, nil

}
