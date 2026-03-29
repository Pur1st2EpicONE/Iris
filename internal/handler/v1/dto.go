package v1

// ShortenLinkDTO represents the payload for creating a shortened link.
// OriginalURL is required, while Alias is optional.
type ShortenLinkDTO struct {
	OriginalURL string `json:"original_url" binding:"required"`
	Alias       string `json:"alias,omitempty"`
}
