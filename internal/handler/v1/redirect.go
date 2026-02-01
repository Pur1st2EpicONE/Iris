package v1

import (
	"Iris/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

const timeout = 5 * time.Second

func (h *Handler) Redirect(c *ginext.Context) {

	shortURL := c.Param("short_url")

	originalURL, err := h.service.GetOriginalURL(c.Request.Context(), models.ShortLink{ShortURL: shortURL})
	if err != nil {
		respondError(c, err)
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		h.service.SaveVisit(ctx, shortURL, c.Request.UserAgent())
	}()

	c.Redirect(http.StatusFound, originalURL)

}
