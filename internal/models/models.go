package models

type Link struct {
	OriginalURL string
	Alias       string
}

type ShortLink struct {
	ShortURL string
}

type VisitStats struct {
	Count       int            `json:"count"`
	ByDay       map[string]int `json:"by_day"`
	ByMonth     map[string]int `json:"by_month"`
	ByUserAgent map[string]int `json:"by_user_agent"`
}
