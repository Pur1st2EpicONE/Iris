package impl

import (
	"Iris/internal/models"
	"errors"
	"net/url"
)

func validateLink(link models.Link) error {

	if err := validateOriginalURL(link.OriginalURL); err != nil {
		return err
	}

	if link.DesiredLength > 0 {
		if err := validateLength(link.DesiredLength); err != nil {
			return err
		}
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
		return errors.New("original_url is not a valid URL")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("original_url must use http or https scheme")
	}

	if u.Host == "" {
		return errors.New("original_url must contain host")
	}

	return nil

}

func validateLength(length uint8) error {
	if length > 100 {
		return errors.New("desired link length exceeds limits")
	}
	return nil
}

func validateAlias(alias string) error {

	if len(alias) > 32 {
		return errors.New("custom_alias length must be less than 32 runes")
	}

	for _, r := range alias {
		if !isAllowedAliasChar(r) {
			return errors.New("custom_alias contains invalid characters")
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
