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

	"github.com/stretchr/testify/require"
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

	logger, _ := logger.NewLogger(config.Logger{Debug: true})

	db, err := dbpg.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode), nil, &dbpg.Options{})
	if err != nil {
		logger.LogFatal("postgres-test — failed to connect to test DB", err, "layer", "repository.postgres-test")
	}

	testStorage = postgres.NewStorage(logger, cfg, db)

	exitCode := m.Run()
	testStorage.Close()

	os.Exit(exitCode)

}

func TestSaveOriginal(t *testing.T) {

	ctx := context.Background()
	originalURL := fmt.Sprintf("https://example.com/%d", time.Now().UnixNano())

	id, err := testStorage.SaveOriginal(ctx, originalURL)
	require.NoError(t, err)
	require.Greater(t, id, int64(0))

}

func TestSaveShort(t *testing.T) {

	ctx := context.Background()
	originalURL := fmt.Sprintf("https://example.com/short-%d", time.Now().UnixNano())

	id, err := testStorage.SaveOriginal(ctx, originalURL)
	require.NoError(t, err)

	shortLink := fmt.Sprintf("short-%d", id)
	err = testStorage.SaveShort(ctx, id, shortLink)
	require.NoError(t, err)

	got, err := testStorage.GetOriginalURL(ctx, shortLink)
	require.NoError(t, err)
	require.Equal(t, originalURL, got)

}

func TestSaveWithAlias(t *testing.T) {

	ctx := context.Background()
	alias := fmt.Sprintf("alias-%d", time.Now().UnixNano())
	link := models.Link{
		OriginalURL: fmt.Sprintf("https://example.com/%s", alias),
		Alias:       alias,
	}

	require.NoError(t, testStorage.SaveWithAlias(ctx, link))
	err := testStorage.SaveWithAlias(ctx, link)
	require.ErrorIs(t, err, errs.ErrAliasExists)

}

func TestGetOriginalURL(t *testing.T) {

	ctx := context.Background()
	alias := fmt.Sprintf("get-%d", time.Now().UnixNano())
	link := models.Link{
		OriginalURL: fmt.Sprintf("https://example.com/%s", alias),
		Alias:       alias,
	}

	require.NoError(t, testStorage.SaveWithAlias(ctx, link))

	url, err := testStorage.GetOriginalURL(ctx, link.Alias)
	require.NoError(t, err)
	require.Equal(t, link.OriginalURL, url)

	_, err = testStorage.GetOriginalURL(ctx, "non-existent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows in result set")

}

func TestSaveVisit(t *testing.T) {

	ctx := context.Background()
	alias := fmt.Sprintf("visit-%d", time.Now().UnixNano())
	link := models.Link{
		OriginalURL: fmt.Sprintf("https://example.com/%s", alias),
		Alias:       alias,
	}

	require.NoError(t, testStorage.SaveWithAlias(ctx, link))

	err := testStorage.SaveVisit(ctx, link.Alias, "GoTest-Agent")
	require.NoError(t, err)

	err = testStorage.SaveVisit(ctx, "non-existent", "GoTest-Agent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows in result set")

}

func countByUA(stats *models.VisitStats, ua string) int {
	count := 0
	for _, entry := range stats.Data {
		if entry.UserAgent == ua {
			count += entry.Count
		}
	}
	return count
}

func TestGetAnalytics(t *testing.T) {

	ctx := context.Background()
	alias := fmt.Sprintf("analytics-%d", time.Now().UnixNano())
	link := models.Link{
		OriginalURL: fmt.Sprintf("https://example.com/%s", alias),
		Alias:       alias,
	}

	require.NoError(t, testStorage.SaveWithAlias(ctx, link))

	require.NoError(t, testStorage.SaveVisit(ctx, link.Alias, "UA1"))
	time.Sleep(1 * time.Millisecond)
	require.NoError(t, testStorage.SaveVisit(ctx, link.Alias, "UA2"))
	time.Sleep(1 * time.Millisecond)
	require.NoError(t, testStorage.SaveVisit(ctx, link.Alias, "UA1"))

	stats, err := testStorage.GetAnalytics(ctx, "", link.Alias)
	require.NoError(t, err)
	require.Equal(t, 3, stats.Count)
	require.Equal(t, 2, countByUA(stats, "UA1"))
	require.Equal(t, 1, countByUA(stats, "UA2"))

	statsUA, err := testStorage.GetAnalytics(ctx, "user_agent", link.Alias)
	require.NoError(t, err)
	require.Equal(t, 3, statsUA.Count)

	countUA1 := 0
	countUA2 := 0

	for _, e := range statsUA.Data {
		switch e.UserAgent {
		case "UA1":
			countUA1 += e.Count
		case "UA2":
			countUA2 += e.Count
		}
	}

	require.Equal(t, 2, countUA1)
	require.Equal(t, 1, countUA2)

	statsDay, err := testStorage.GetAnalytics(ctx, "day", link.Alias)
	require.NoError(t, err)
	require.Equal(t, 3, statsDay.Count)
	require.Len(t, statsDay.Data, 3)

	for _, e := range statsDay.Data {
		require.NotEmpty(t, e.Key)
	}

	statsMonth, err := testStorage.GetAnalytics(ctx, "month", link.Alias)
	require.NoError(t, err)
	require.Equal(t, 3, statsMonth.Count)
	require.Len(t, statsMonth.Data, 3)

	for _, e := range statsMonth.Data {
		require.NotEmpty(t, e.Key)
	}

	_, err = testStorage.GetAnalytics(ctx, "invalid", link.Alias)
	require.ErrorIs(t, err, errs.ErrInvalidGroupBy)

	_, err = testStorage.GetAnalytics(ctx, "", "non-existent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows in result set")

}
