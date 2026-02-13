package v1

import (
	"Iris/internal/errs"
	"Iris/internal/models"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) Shorten(c *ginext.Context) {

	var request ShortenLinkDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	shortLink, err := h.service.ShortenLink(c.Request.Context(),
		models.Link{OriginalURL: request.OriginalURL, Alias: request.Alias})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, shortLink)

}
