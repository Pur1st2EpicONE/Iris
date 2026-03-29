package v1

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	mockService "Iris/internal/service/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/ginext"
)

func setupRouter(handler *Handler) *ginext.Engine {

	r := ginext.New("")

	v1 := r.Group("/v1")
	{
		v1.POST("/shorten", handler.Shorten)
		v1.GET("/:short_url", handler.Redirect)
		v1.GET("/:short_url/analytics", handler.GetAnalytics)
	}

	return r

}

func TestHandler_Shorten(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mockService.NewMockService(controller)

	h := NewHandler(mockService)
	router := setupRouter(h)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/shorten", bytes.NewBufferString(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns error", func(t *testing.T) {
		body := ShortenLinkDTO{OriginalURL: "http://test.com"}
		b, _ := json.Marshal(body)
		mockService.EXPECT().ShortenLink(gomock.Any(), models.Link{OriginalURL: body.OriginalURL}).Return("", errs.ErrAliasExists)
		req := httptest.NewRequest(http.MethodPost, "/v1/shorten", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		body := ShortenLinkDTO{OriginalURL: "http://test.com"}
		b, _ := json.Marshal(body)
		mockService.EXPECT().ShortenLink(gomock.Any(), models.Link{OriginalURL: body.OriginalURL}).Return("abc123", nil)
		req := httptest.NewRequest(http.MethodPost, "/v1/shorten", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "abc123")
	})
}

func TestHandler_Redirect(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mockService.NewMockService(controller)

	h := NewHandler(mockService)
	router := setupRouter(h)

	t.Run("link not found", func(t *testing.T) {
		mockService.EXPECT().GetOriginalURL(gomock.Any(), models.ShortLink{ShortURL: "abc"}).Return("", errs.ErrLinkNotFound)
		req := httptest.NewRequest(http.MethodGet, "/v1/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("success redirect", func(t *testing.T) {

		var wg sync.WaitGroup
		wg.Add(1)

		mockService.EXPECT().GetOriginalURL(gomock.Any(), models.ShortLink{ShortURL: "abc"}).Return("https://google.com", nil)
		mockService.EXPECT().SaveVisit(gomock.Any(), "abc", gomock.Any()).DoAndReturn(func(ctx context.Context, shortURL, userAgent string) error {
			defer wg.Done()
			return nil
		})

		req := httptest.NewRequest(http.MethodGet, "/v1/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusFound, w.Code)
		require.Equal(t, "https://google.com", w.Header().Get("Location"))

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Error("Timeout waiting for goroutine")
		}
	})

}

func TestHandler_GetAnalytics(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mockService.NewMockService(controller)

	h := NewHandler(mockService)
	router := setupRouter(h)

	stats := &models.VisitStats{
		Count: 12,
		Data: []models.VisitEntry{
			{Key: "2026-02-01", UserAgent: "Chrome", Time: "2026-02-01", Count: 5},
			{Key: "2026-02-02", UserAgent: "Safari", Time: "2026-02-02", Count: 7},
		},
	}

	t.Run("invalid group_by", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/abc/analytics?group_by=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), errs.ErrInvalidGroupBy.Error())
	})

	t.Run("link not found", func(t *testing.T) {
		mockService.EXPECT().GetAnalytics(gomock.Any(), "", "abc").Return(nil, errs.ErrLinkNotFound)
		req := httptest.NewRequest(http.MethodGet, "/v1/abc/analytics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		mockService.EXPECT().GetAnalytics(gomock.Any(), "", "abc").Return(nil, errors.New("db down"))
		req := httptest.NewRequest(http.MethodGet, "/v1/abc/analytics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success default group_by", func(t *testing.T) {
		mockService.EXPECT().GetAnalytics(gomock.Any(), "", "abc").Return(stats, nil)
		req := httptest.NewRequest(http.MethodGet, "/v1/abc/analytics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Result models.VisitStats `json:"result"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, 12, resp.Result.Count)
		require.Len(t, resp.Result.Data, 2)
		require.Equal(t, "Chrome", resp.Result.Data[0].UserAgent)
		require.Equal(t, "Safari", resp.Result.Data[1].UserAgent)
	})

	t.Run("success with group_by=day", func(t *testing.T) {
		mockService.EXPECT().GetAnalytics(gomock.Any(), "day", "abc").Return(stats, nil)
		req := httptest.NewRequest(http.MethodGet, "/v1/abc/analytics?group_by=day", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Result models.VisitStats `json:"result"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, 12, resp.Result.Count)
		require.Len(t, resp.Result.Data, 2)
	})

}
