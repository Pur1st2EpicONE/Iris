package v1

import (
	"Iris/internal/errs"
	"errors"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

func respondOK(c *ginext.Context, response any) {
	c.JSON(http.StatusOK, ginext.H{"result": response})
}

func respondError(c *ginext.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, ginext.H{"error": msg})
	}
}

func mapErrorToStatus(err error) (int, string) {

	switch {
	case errors.Is(err, errs.ErrInvalidJSON),
		errors.Is(err, errs.ErrInvalidOriginalURL),
		errors.Is(err, errs.ErrOriginalURLScheme),
		errors.Is(err, errs.ErrOriginalURLHost),
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
