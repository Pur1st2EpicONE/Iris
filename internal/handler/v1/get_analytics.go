package v1

import (
	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) GetAnalytics(c *ginext.Context) {
	data, err := h.service.GetAnalytics(c.Request.Context(), c.Param("short_url"))
	if err != nil {
		respondError(c, err)
	} else {
		respondOK(c, data)
	}
}
