package impl

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"net/url"
)

// validateLink checks that a Link's original URL and alias (if present) are valid.
// Returns an error if any validation rule fails.
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

// validateOriginalURL verifies that the URL is parseable, has a supported scheme (http/https),
// and contains a host. Returns appropriate errors from errs package.
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

// validateAlias checks that the alias is at most 32 characters long
// and contains only allowed characters (letters, digits, '-' or '_').
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

// isAllowedAliasChar returns true if the rune is a valid character in a short link alias.
func isAllowedAliasChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_'
}
