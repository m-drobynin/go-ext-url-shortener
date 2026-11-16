package model

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortenBatchRequestRow struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenBatchResponseRow struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UrlsPairRow struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
