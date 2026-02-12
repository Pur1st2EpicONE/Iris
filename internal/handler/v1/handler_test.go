package v1

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	mockService "Iris/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockService(ctrl)

	h := &Handler{service: mockService}
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockService(ctrl)

	h := &Handler{service: mockService}
	router := setupRouter(h)

	t.Run("link not found", func(t *testing.T) {
		mockService.EXPECT().GetOriginalURL(gomock.Any(), models.ShortLink{ShortURL: "abc"}).Return("", errs.ErrLinkNotFound)
		req := httptest.NewRequest(http.MethodGet, "/v1/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("success redirect", func(t *testing.T) {
		mockService.EXPECT().GetOriginalURL(gomock.Any(), models.ShortLink{ShortURL: "abc"}).Return("https://google.com", nil)
		mockService.EXPECT().SaveVisit(gomock.Any(), "abc", gomock.Any()).AnyTimes()
		req := httptest.NewRequest(http.MethodGet, "/v1/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusFound, w.Code)
		require.Equal(t, "https://google.com", w.Header().Get("Location"))
	})

}

func TestHandler_GetAnalytics(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockService(ctrl)

	h := &Handler{service: mockService}
	router := setupRouter(h)

	stats := &models.VisitStats{
		Count: 12,
		ByDay: map[string]int{
			"2026-02-01": 5,
			"2026-02-02": 7,
		},
		ByMonth: map[string]int{
			"2026-02": 12,
		},
		ByUserAgent: map[string]int{
			"Chrome": 8,
			"Safari": 4,
		},
	}

	t.Run("link not found", func(t *testing.T) {
		mockService.EXPECT().GetAnalytics(gomock.Any(), "abc").Return(nil, errs.ErrLinkNotFound)
		req := httptest.NewRequest(http.MethodGet, "/v1/abc/analytics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		mockService.EXPECT().GetAnalytics(gomock.Any(), "abc").Return(nil, errors.New("db down"))
		req := httptest.NewRequest(http.MethodGet, "/v1/abc/analytics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		mockService.EXPECT().GetAnalytics(gomock.Any(), "abc").Return(stats, nil)
		req := httptest.NewRequest(http.MethodGet, "/v1/abc/analytics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, `"count":12`)
		require.Contains(t, body, `"by_day"`)
		require.Contains(t, body, `"by_month"`)
		require.Contains(t, body, `"by_user_agent"`)
	})

}
