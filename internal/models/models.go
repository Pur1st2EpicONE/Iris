package models

type Link struct {
	OriginalURL string
	Alias       string
}

type ShortLink struct {
	ShortURL string
}

type VisitEntry struct {
	Key       string
	UserAgent string
	Time      string
	Count     int
}

type VisitStats struct {
	Count int
	Data  []VisitEntry
}
