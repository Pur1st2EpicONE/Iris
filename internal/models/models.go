package models

type Link struct {
	OriginalURL string
	Alias       string
}

type ShortLink struct {
	ShortURL string
}

type VisitEntry struct {
	Key       string `json:"key"`
	UserAgent string `json:"user_agent"`
	Time      string `json:"time"`
	Count     int    `json:"count"`
}

type VisitStats struct {
	Count int          `json:"count"`
	Data  []VisitEntry `json:"data"`
}
