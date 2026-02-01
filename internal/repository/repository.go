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

type Storage interface {
	SaveOriginal(ctx context.Context, originalURL string) (int64, error)
	SaveWithAlias(ctx context.Context, link models.Link) error
	SaveShort(ctx context.Context, id int64, shortURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	SaveVisit(ctx context.Context, shortURL string, userAgent string) error
	GetAnalytics(ctx context.Context, shortURL string) (*models.VisitStats, error)
	Close()
}

func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) Storage {
	return postgres.NewStorage(logger, config, db)
}

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
