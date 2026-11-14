package model

type ShortenedURLRow struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserId      string `json:"user_id"`
}

type ValuesPair struct {
	Key   string
	Value string
}

type User struct {
	UserID string
}
