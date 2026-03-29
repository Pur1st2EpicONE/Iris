package v1

import (
	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) GetAnalytics(c *ginext.Context) {

	groupBy, err := parseQuery(c)
	if err != nil {
		respondError(c, err)
		return
	}

	analytics, err := h.service.GetAnalytics(c.Request.Context(), groupBy, c.Param("short_url"))
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, analytics)

}
