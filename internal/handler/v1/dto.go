package v1

type ShortenLinkDTO struct {
	OriginalURL string `json:"original_url" binding:"required"`
	Alias       string `json:"alias,omitempty"`
}
