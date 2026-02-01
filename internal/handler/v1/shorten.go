package v1

import (
	"Iris/internal/errs"
	"Iris/internal/models"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) Shorten(c *ginext.Context) {

	var r ShortenLinkDTO

	if err := c.ShouldBindJSON(&r); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	shortLink, err := h.service.ShortenLink(c.Request.Context(), models.Link{OriginalURL: r.OriginalURL, Alias: r.Alias})
	if err != nil {
		respondError(c, err)
	} else {
		respondOK(c, shortLink)
	}

}
