package v1

type ShortenLinkV1 struct {
	OriginalURL   string `json:"original_url" binding:"required"`
	DesiredLength uint8  `json:"desired_length,omitempty"`
	Alias         string `json:"custom_alias,omitempty"`
}
