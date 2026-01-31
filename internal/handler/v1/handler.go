package v1

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"Iris/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ShortenLink(c *ginext.Context) {

	var request ShortenLinkV1

	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	link := models.Link{OriginalURL: request.OriginalURL, DesiredLength: request.DesiredLength, Alias: request.Alias}

	id, err := h.service.ShortenLink(c.Request.Context(), link)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, id)

}
