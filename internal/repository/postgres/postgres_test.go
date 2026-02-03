package postgres_test

import (
	"Iris/internal/config"
	"Iris/internal/errs"
	"Iris/internal/logger"
	"Iris/internal/models"
	"Iris/internal/repository/postgres"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	wbf "github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/dbpg"
)

var testStorage *postgres.Storage

func TestMain(m *testing.M) {

	if err := wbf.New().LoadEnvFiles("../../../.env"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := config.Storage{
		Host:     "postgres-test",
		Port:     "5432",
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   "iris_test",
		SSLMode:  "disable",
		QueryRetryStrategy: config.RetryStrategy{
			Attempts: 3,
			Delay:    100 * time.Millisecond,
			Backoff:  1.5,
		},
	}

	log, _ := logger.NewLogger(config.Logger{Debug: true})

	db, err := dbpg.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode), nil, &dbpg.Options{})
	if err != nil {
		log.LogFatal("postgres-test — failed to connect to test DB", err, "layer", "repository.postgres-test")
	}

	testStorage = postgres.NewStorage(log, cfg, db)

	exitCode := m.Run()
	testStorage.Close()

	os.Exit(exitCode)

}

func TestSaveOriginal(t *testing.T) {

	ctx := context.Background()
	originalURL := "https://example.com"

	id, err := testStorage.SaveOriginal(ctx, originalURL)
	if err != nil {
		t.Fatalf("SaveOriginal failed: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected id > 0, got %d", id)
	}

}

func TestSaveShort(t *testing.T) {

	ctx := context.Background()
	originalURL := "https://example.com/short"

	id, _ := testStorage.SaveOriginal(ctx, originalURL)
	shortLink := fmt.Sprintf("short-%d", time.Now().UnixNano())

	err := testStorage.SaveShort(ctx, id, shortLink)
	if err != nil {
		t.Fatalf("SaveShort failed: %v", err)
	}

}

func TestSaveWithAlias(t *testing.T) {

	ctx := context.Background()
	link := models.Link{
		OriginalURL: "https://example.com/alias",
		Alias:       fmt.Sprintf("alias-%d", time.Now().UnixNano()),
	}

	if err := testStorage.SaveWithAlias(ctx, link); err != nil {
		t.Fatalf("SaveWithAlias failed: %v", err)
	}

	err := testStorage.SaveWithAlias(ctx, link)
	if err != errs.ErrAliasExists {
		t.Fatalf("expected ErrAliasExists, got %v", err)
	}

}

func TestGetOriginalURL(t *testing.T) {

	ctx := context.Background()
	link := models.Link{
		OriginalURL: "https://example.com/get",
		Alias:       fmt.Sprintf("get-%d", time.Now().UnixNano()),
	}
	_ = testStorage.SaveWithAlias(ctx, link)

	url, err := testStorage.GetOriginalURL(ctx, link.Alias)
	if err != nil {
		t.Fatalf("GetOriginalURL failed: %v", err)
	}

	if url != link.OriginalURL {
		t.Fatalf("expected %s, got %s", link.OriginalURL, url)
	}

	_, err = testStorage.GetOriginalURL(ctx, "non-existent")
	if err != errs.ErrLinkNotFound {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}

}

func TestSaveVisit(t *testing.T) {

	ctx := context.Background()
	link := models.Link{
		OriginalURL: "https://example.com/visit",
		Alias:       fmt.Sprintf("visit-%d", time.Now().UnixNano()),
	}
	_ = testStorage.SaveWithAlias(ctx, link)

	err := testStorage.SaveVisit(ctx, link.Alias, "GoTest-Agent")
	if err != nil {
		t.Fatalf("SaveVisit failed: %v", err)
	}

	err = testStorage.SaveVisit(ctx, "non-existent", "GoTest-Agent")
	if err != errs.ErrLinkNotFound {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}

}

func TestGetAnalytics(t *testing.T) {

	ctx := context.Background()
	link := models.Link{
		OriginalURL: "https://example.com/analytics",
		Alias:       fmt.Sprintf("analytics-%d", time.Now().UnixNano()),
	}
	_ = testStorage.SaveWithAlias(ctx, link)

	_ = testStorage.SaveVisit(ctx, link.Alias, "UA1")
	_ = testStorage.SaveVisit(ctx, link.Alias, "UA2")
	_ = testStorage.SaveVisit(ctx, link.Alias, "UA1")

	stats, err := testStorage.GetAnalytics(ctx, link.Alias)
	if err != nil {
		t.Fatalf("GetAnalytics failed: %v", err)
	}

	if stats.Count != 3 {
		t.Fatalf("expected 3 visits, got %d", stats.Count)
	}

	if stats.ByUserAgent["UA1"] != 2 || stats.ByUserAgent["UA2"] != 1 {
		t.Fatalf("user agent counts mismatch: %v", stats.ByUserAgent)
	}

	_, err = testStorage.GetAnalytics(ctx, "non-existent")
	if err != errs.ErrLinkNotFound {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}

}
