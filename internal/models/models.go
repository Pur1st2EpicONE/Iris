// Package models defines the core data structures for the Iris application.
// It includes models for links, shortened links, and visit statistics.
package models

// Link represents an original URL that can optionally have a custom alias.
type Link struct {
	OriginalURL string // OriginalURL is the full URL to be shortened
	Alias       string // Alias is an optional custom short identifier
}

// ShortLink represents a generated short URL.
type ShortLink struct {
	ShortURL string // ShortURL is the shortened version of the original link
}

// VisitEntry represents a single aggregated visit entry for analytics purposes.
type VisitEntry struct {
	Key       string `json:"key"`        // Key is the grouping key (e.g., date, user agent)
	UserAgent string `json:"user_agent"` // UserAgent of the visitor
	Time      string `json:"time"`       // Time of the visit
	Count     int    `json:"count"`      // Number of visits for this key
}

// VisitStats represents summarized visit statistics for a short link.
type VisitStats struct {
	Count int          `json:"count"` // Count is the total number of visits
	Data  []VisitEntry `json:"data"`  // Data contains individual visit entries
}
