package model

import (
	"sync"
	"uuid"
)

type Database struct {
	mu   sync.Mutex
	rows []ShortenedURLRow
}

type ShortenedURLRow struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewDatabase() *Database {
	var database = &Database{
		rows: []ShortenedURLRow{},
	}

	return database
}

func (db *Database) TryPut(key string, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.get(key)

	if ok {
		return ErrDBConflict
	}

	newRow := ShortenedURLRow{
		ShortURL:    key,
		OriginalURL: value,
		UUID:        uuid.NewV4().String(),
	}

	db.rows = append(db.rows, newRow)

	return nil
}

func (db *Database) Get(key string) (*string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, ok := db.get(key)

	if ok {
		return res, nil
	}

	return nil, ErrDBNotFound
}

func (db *Database) get(key string) (*string, bool) {
	for _, v := range db.rows {
		if v.ShortURL == key {
			return &v.OriginalURL, true
		}
	}

	return nil, false
}
