package impl

import (
	mockCache "Iris/internal/cache/mocks"
	"Iris/internal/errs"
	mockLogger "Iris/internal/logger/mocks"
	"Iris/internal/models"
	mockStorage "Iris/internal/repository/mocks"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService_GetAnalytics(t *testing.T) {

	ctx := context.Background()
	shortURL := "qweqwe"

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockCache := mockCache.NewMockCache(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	service := NewService(mockLogger, mockCache, mockStorage)

	stats := &models.VisitStats{
		Count:       3,
		ByUserAgent: map[string]int{"UA1": 2, "UA2": 1},
	}

	t.Run("success", func(t *testing.T) {
		mockStorage.EXPECT().GetAnalytics(ctx, shortURL).Return(stats, nil)
		res, err := service.GetAnalytics(ctx, shortURL)
		require.NoError(t, err)
		require.Equal(t, stats, res)
	})

	t.Run("link not found", func(t *testing.T) {
		mockStorage.EXPECT().GetAnalytics(ctx, shortURL).Return(nil, errs.ErrLinkNotFound)
		res, err := service.GetAnalytics(ctx, shortURL)
		require.ErrorIs(t, err, errs.ErrLinkNotFound)
		require.Nil(t, res)
	})

	t.Run("storage error", func(t *testing.T) {
		dbErr := errors.New("db down")
		mockStorage.EXPECT().GetAnalytics(ctx, shortURL).Return(nil, dbErr)
		mockLogger.EXPECT().LogError("service — failed to get analytics", dbErr, "short link", shortURL, "layer", "service.impl")
		res, err := service.GetAnalytics(ctx, shortURL)
		require.Error(t, err)
		require.Nil(t, res)
	})

}

func TestService_GetOriginalURL(t *testing.T) {

	ctx := context.Background()
	link := models.ShortLink{ShortURL: "qweqwe"}
	original := "https://example.com"

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockCache := mockCache.NewMockCache(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	service := NewService(mockLogger, mockCache, mockStorage)

	t.Run("found in cache", func(t *testing.T) {
		mockCache.EXPECT().GetLink(ctx, link.ShortURL).Return(original, nil)
		mockLogger.EXPECT().Debug("service — link fetched from cache", "short link", link.ShortURL, "layer", "service.impl")
		res, err := service.GetOriginalURL(ctx, link)
		require.NoError(t, err)
		require.Equal(t, original, res)
	})

	t.Run("cache miss, found in DB, cache set ok", func(t *testing.T) {
		mockCache.EXPECT().GetLink(ctx, link.ShortURL).Return("", errors.New("cache miss"))
		mockStorage.EXPECT().GetOriginalURL(ctx, link.ShortURL).Return(original, nil)
		mockCache.EXPECT().SetLink(ctx, link.ShortURL, original).Return(nil)
		mockLogger.EXPECT().Debug("service — link fetched from DB", "short link", link.ShortURL, "layer", "service.impl")
		res, err := service.GetOriginalURL(ctx, link)
		require.NoError(t, err)
		require.Equal(t, original, res)
	})

	t.Run("cache miss, DB not found", func(t *testing.T) {
		mockCache.EXPECT().GetLink(ctx, link.ShortURL).Return("", errors.New("cache miss"))
		mockStorage.EXPECT().GetOriginalURL(ctx, link.ShortURL).Return("", errs.ErrLinkNotFound)
		res, err := service.GetOriginalURL(ctx, link)
		require.ErrorIs(t, err, errs.ErrLinkNotFound)
		require.Empty(t, res)
	})

	t.Run("cache miss, DB error", func(t *testing.T) {
		dbErr := errors.New("db down")
		mockCache.EXPECT().GetLink(ctx, link.ShortURL).Return("", errors.New("cache miss"))
		mockStorage.EXPECT().GetOriginalURL(ctx, link.ShortURL).Return("", dbErr)
		mockLogger.EXPECT().LogError("service — failed to get original url from DB", dbErr, "short link", link.ShortURL, "layer", "service.impl")
		res, err := service.GetOriginalURL(ctx, link)
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("cache set fails but request succeeds", func(t *testing.T) {
		mockCache.EXPECT().GetLink(ctx, link.ShortURL).Return("", errors.New("cache miss"))
		mockStorage.EXPECT().GetOriginalURL(ctx, link.ShortURL).Return(original, nil)
		mockCache.EXPECT().SetLink(ctx, link.ShortURL, original).Return(errors.New("cache down"))
		mockLogger.EXPECT().LogError("service — failed to save link in cache", gomock.Any(), "short link", link.ShortURL, "layer", "service.impl")
		mockLogger.EXPECT().Debug("service — link fetched from DB", "short link", link.ShortURL, "layer", "service.impl")
		res, err := service.GetOriginalURL(ctx, link)
		require.NoError(t, err)
		require.Equal(t, original, res)
	})

}

func TestService_SaveVisit(t *testing.T) {

	ctx := context.Background()
	shortURL := "qweqwe"
	ua := "asd"

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockCache := mockCache.NewMockCache(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	service := NewService(mockLogger, mockCache, mockStorage)

	t.Run("success", func(t *testing.T) {
		mockStorage.EXPECT().SaveVisit(ctx, shortURL, ua).Return(nil)
		service.SaveVisit(ctx, shortURL, ua)
	})

	t.Run("storage error", func(t *testing.T) {
		err := errors.New("db down")
		mockStorage.EXPECT().SaveVisit(ctx, shortURL, ua).Return(err)
		mockLogger.EXPECT().LogError("service — failed to save visit", err, "short link", shortURL, "user agent", ua, "layer", "service.impl")
		service.SaveVisit(ctx, shortURL, ua)
	})

}

func TestService_ShortenLink(t *testing.T) {

	ctx := context.Background()
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockCache := mockCache.NewMockCache(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	service := NewService(mockLogger, mockCache, mockStorage)

	t.Run("invalid link", func(t *testing.T) {
		_, err := service.ShortenLink(ctx, models.Link{})
		require.Error(t, err)
	})

	t.Run("save with alias success", func(t *testing.T) {
		link := models.Link{OriginalURL: "https://example.com", Alias: "custom"}
		mockStorage.EXPECT().SaveWithAlias(ctx, link).Return(nil)
		res, err := service.ShortenLink(ctx, link)
		require.NoError(t, err)
		require.Equal(t, "custom", res)
	})

	t.Run("alias already exists", func(t *testing.T) {
		link := models.Link{OriginalURL: "https://example.com", Alias: "custom"}
		mockStorage.EXPECT().SaveWithAlias(ctx, link).Return(errs.ErrAliasExists)
		_, err := service.ShortenLink(ctx, link)
		require.ErrorIs(t, err, errs.ErrAliasExists)
	})

	t.Run("auto generate short link", func(t *testing.T) {
		link := models.Link{OriginalURL: "https://example.com"}
		mockStorage.EXPECT().SaveOriginal(ctx, link.OriginalURL).Return(int64(1), nil)
		mockStorage.EXPECT().SaveShort(ctx, int64(1), gomock.Any()).Return(nil)
		res, err := service.ShortenLink(ctx, link)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})

	t.Run("invalid original URL scheme", func(t *testing.T) {
		link := models.Link{OriginalURL: "ftp://example.com"}
		_, err := service.ShortenLink(ctx, link)
		require.ErrorIs(t, err, errs.ErrOriginalURLScheme)
	})

	t.Run("invalid original URL host missing", func(t *testing.T) {
		link := models.Link{OriginalURL: "http://"}
		_, err := service.ShortenLink(ctx, link)
		require.ErrorIs(t, err, errs.ErrOriginalURLHost)
	})

	t.Run("alias too long", func(t *testing.T) {
		alias := strings.Repeat("a", 33)
		link := models.Link{OriginalURL: "https://example.com", Alias: alias}
		_, err := service.ShortenLink(ctx, link)
		require.ErrorIs(t, err, errs.ErrAliasTooLong)
	})

	t.Run("alias invalid chars", func(t *testing.T) {
		link := models.Link{OriginalURL: "https://example.com", Alias: "abc!"}
		_, err := service.ShortenLink(ctx, link)
		require.ErrorIs(t, err, errs.ErrAliasInvalidChars)
	})

	t.Run("SaveOriginal error", func(t *testing.T) {
		link := models.Link{OriginalURL: "https://example.com"}
		dbErr := errors.New("db save original error")
		mockStorage.EXPECT().SaveOriginal(ctx, link.OriginalURL).Return(int64(0), dbErr)
		mockLogger.EXPECT().LogError("service — failed to save original link", gomock.Any(),
			"link", link.OriginalURL, "layer", "service.impl")
		_, err := service.ShortenLink(ctx, link)
		require.ErrorIs(t, err, dbErr)
	})

	t.Run("SaveShort error", func(t *testing.T) {
		link := models.Link{OriginalURL: "https://example.com"}
		id := int64(1)
		shortLink := encode(id + offset)
		dbErr := errors.New("db save short error")
		mockStorage.EXPECT().SaveOriginal(ctx, link.OriginalURL).Return(id, nil)
		mockStorage.EXPECT().SaveShort(ctx, id, shortLink).Return(dbErr)
		mockLogger.EXPECT().LogError("service — failed to save short link", gomock.Any(),
			"link", link.OriginalURL, "short link", shortLink, "layer", "service.impl")
		_, err := service.ShortenLink(ctx, link)
		require.ErrorIs(t, err, dbErr)
	})

	t.Run("saveWithAlias unexpected error", func(t *testing.T) {
		link := models.Link{OriginalURL: "https://example.com", Alias: "custom"}
		unexpectedErr := errors.New("unexpected error")
		mockStorage.EXPECT().SaveWithAlias(ctx, link).Return(unexpectedErr)
		mockLogger.EXPECT().LogError("service — failed to save link with alias", gomock.Any(),
			"link", link.OriginalURL, "alias", link.Alias, "layer", "service.impl")
		_, err := service.ShortenLink(ctx, link)
		require.ErrorIs(t, err, unexpectedErr)
	})

}
