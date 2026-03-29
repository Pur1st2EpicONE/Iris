// Package repository provides an abstraction layer over the database for
// managing URL shortening data, including links, short URLs, visits, and analytics.
package repository

import (
	"Iris/internal/config"
	"Iris/internal/logger"
	"Iris/internal/models"
	"Iris/internal/repository/postgres"
	"context"
	"fmt"

	"github.com/wb-go/wbf/dbpg"
)

// Storage defines the interface for interacting with the URL shortener storage.
// It abstracts database operations such as saving original URLs, generating short links,
// recording visits, retrieving original URLs, fetching analytics, and closing the storage.
type Storage interface {
	SaveOriginal(ctx context.Context, originalURL string) (int64, error)                           // SaveOriginal inserts the original URL into the database and returns its ID.
	SaveWithAlias(ctx context.Context, link models.Link) error                                     // SaveWithAlias saves a URL with a custom alias.
	SaveShort(ctx context.Context, id int64, shortURL string) error                                // SaveShort updates the record with the generated short URL for a given ID.
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)                           // GetOriginalURL retrieves the original URL corresponding to the given short URL.
	SaveVisit(ctx context.Context, shortURL string, userAgent string) error                        // SaveVisit records a visit for the given short URL and user agent.
	GetAnalytics(ctx context.Context, groupBy string, shortURL string) (*models.VisitStats, error) // GetAnalytics fetches visit statistics for a given short URL, optionally grouped by day, month, or user agent.
	Close()                                                                                        // Close closes the database connection and cleans up resources.
}

// NewStorage creates a new Storage instance backed by Postgres.
func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) Storage {
	return postgres.NewStorage(logger, config, db)
}

// ConnectDB establishes a connection to the Postgres database using the provided configuration.
// It returns a dbpg.DB instance ready for queries. It performs a ping to verify connectivity.
func ConnectDB(config config.Storage) (*dbpg.DB, error) {

	options := &dbpg.Options{
		MaxOpenConns:    config.MaxOpenConns,
		MaxIdleConns:    config.MaxIdleConns,
		ConnMaxLifetime: config.ConnMaxLifetime,
	}

	db, err := dbpg.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode), nil, options)
	if err != nil {
		return nil, fmt.Errorf("database driver not found or DSN invalid: %w", err)
	}

	if err := db.Master.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil

}
