package errs

import "errors"

var (
	ErrInvalidJSON        = errors.New("invalid JSON format")                                 // invalid JSON format
	ErrLinkNotFound       = errors.New("short link not found")                                // short link not found
	ErrAliasExists        = errors.New("URL with this alias already exists")                  // URL with this alias already exists
	ErrInternal           = errors.New("internal server error")                               // internal server error
	ErrInvalidOriginalURL = errors.New("original_url is not a valid URL")                     // original_url is not a valid URL
	ErrOriginalURLScheme  = errors.New("original_url must use http or https scheme")          // original_url must use http or https scheme
	ErrOriginalURLHost    = errors.New("original_url must contain host")                      // original_url must contain host
	ErrAliasTooLong       = errors.New("custom_alias length must be less than 32 characters") // custom_alias length must be less than 32 characters
	ErrAliasInvalidChars  = errors.New("custom_alias contains invalid characters")            // custom_alias contains invalid characters
)
