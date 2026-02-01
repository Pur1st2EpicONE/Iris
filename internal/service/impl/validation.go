package impl

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"net/url"
)

func validateLink(link models.Link) error {

	if err := validateOriginalURL(link.OriginalURL); err != nil {
		return err
	}

	if link.Alias != "" {
		if err := validateAlias(link.Alias); err != nil {
			return err
		}
	}

	return nil

}

func validateOriginalURL(link string) error {

	u, err := url.ParseRequestURI(link)
	if err != nil {
		return errs.ErrInvalidOriginalURL
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errs.ErrOriginalURLScheme
	}

	if u.Host == "" {
		return errs.ErrOriginalURLHost
	}

	return nil

}

func validateAlias(alias string) error {

	if len(alias) > 32 {
		return errs.ErrAliasTooLong
	}

	for _, r := range alias {
		if !isAllowedAliasChar(r) {
			return errs.ErrAliasInvalidChars
		}
	}

	return nil

}

func isAllowedAliasChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_'
}
