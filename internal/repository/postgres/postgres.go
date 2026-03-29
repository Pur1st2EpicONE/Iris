// Package postgres implements the Postgres-backed repository for the URL shortener.
// It provides methods for saving links, recording visits, and retrieving analytics.
package postgres

import (
	"Iris/internal/config"
	"Iris/internal/logger"

	"github.com/wb-go/wbf/dbpg"
)

// Storage is a Postgres-backed implementation of the repository.Storage interface.
// It provides methods to query and manipulate links and visit data.
type Storage struct {
	db     *dbpg.DB       // db is the Postgres connection
	logger logger.Logger  // logger logs storage-related events
	config config.Storage // config contains storage configuration options
}

// NewStorage creates a new Postgres Storage instance with the given logger, config, and db connection.
func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *Storage {
	return &Storage{db: db, logger: logger, config: config}
}

// Close closes the database connection and logs the result.
func (s *Storage) Close() {
	if err := s.db.Master.Close(); err != nil {
		s.logger.LogError("postgres — failed to close properly", err, "layer", "repository.postgres")
	} else {
		s.logger.LogInfo("postgres — database closed", "layer", "repository.postgres")
	}
}
