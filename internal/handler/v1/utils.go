package v1

import (
	"Iris/internal/errs"
	"errors"
	"net/http"
	"strings"

	"github.com/wb-go/wbf/ginext"
)

// parseQuery validates and returns the optional "group_by" query parameter.
func parseQuery(c *ginext.Context) (string, error) {

	groupBy := strings.ToLower(strings.TrimSpace(c.Query("group_by")))

	switch groupBy {
	case "", "day", "month", "user_agent":
		return groupBy, nil
	default:
		return "", errs.ErrInvalidGroupBy
	}

}

// respondOK sends a successful JSON response with status 200.
func respondOK(c *ginext.Context, response any) {
	c.JSON(http.StatusOK, ginext.H{"result": response})
}

// respondError sends an error response based on the type of error.
func respondError(c *ginext.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, ginext.H{"error": msg})
	}
}

// mapErrorToStatus maps internal errors to HTTP status codes and messages.
func mapErrorToStatus(err error) (int, string) {

	switch {
	case errors.Is(err, errs.ErrInvalidJSON),
		errors.Is(err, errs.ErrInvalidOriginalURL),
		errors.Is(err, errs.ErrOriginalURLScheme),
		errors.Is(err, errs.ErrOriginalURLHost),
		errors.Is(err, errs.ErrInvalidGroupBy),
		errors.Is(err, errs.ErrAliasTooLong),
		errors.Is(err, errs.ErrAliasInvalidChars):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, errs.ErrAliasExists):
		return http.StatusConflict, err.Error()

	case errors.Is(err, errs.ErrLinkNotFound):
		return http.StatusNotFound, err.Error()

	default:
		return http.StatusInternalServerError, errs.ErrInternal.Error()
	}

}
