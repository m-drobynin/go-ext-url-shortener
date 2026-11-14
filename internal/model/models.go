package model

type ShortenedURLRow struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ValuesPair struct {
	Key   string
	Value string
}

type User struct {
	UserID string
}
