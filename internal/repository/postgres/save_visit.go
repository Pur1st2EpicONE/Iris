package postgres

import "context"

func (s *Storage) SaveVisit(ctx context.Context, shortURL string, userAgent string) error {

	var linkID int64
	query := `
	
	SELECT id 
	FROM links 
	WHERE short_link = $1`

	err := s.db.QueryRowContext(ctx, query, shortURL).Scan(&linkID)
	if err != nil {
		return err
	}

	insert := `

    INSERT INTO visits (link_id, user_agent)
    VALUES ($1, $2)`

	_, err = s.db.ExecContext(ctx, insert, linkID, userAgent)
	return err

}
