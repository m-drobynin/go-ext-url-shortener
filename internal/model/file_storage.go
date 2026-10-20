package model

import (
	"encoding/json"
	"os"
	"uuid"
)

type FileStorage struct {
	filename string
	rows     []ShortenedURLRow
}

type ShortenedURLRow struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(filename string) (*FileStorage, error) {
	fs := FileStorage{
		filename: filename,
		rows:     []ShortenedURLRow{},
	}

	if err := fs.read(); err != nil {
		return nil, err
	}

	return &fs, nil
}

func (fs *FileStorage) TryPut(key string, value string) error {
	newRow := ShortenedURLRow{
		ShortURL:    key,
		OriginalURL: value,
		UUID:        uuid.NewV4().String(),
	}

	fs.rows = append(fs.rows, newRow)

	return fs.save()
}

func (fs *FileStorage) Get(key string) (*string, error) {
	for _, v := range fs.rows {
		if v.ShortURL == key {
			return &v.OriginalURL, nil
		}
	}

	return nil, nil
}

func (fs *FileStorage) save() error {
	data, err := json.Marshal(fs.rows)

	if err != nil {
		return err
	}

	return os.WriteFile(fs.filename, data, 0666)
}

func (fs *FileStorage) read() error {
	data, err := os.ReadFile(fs.filename)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &fs.rows); err != nil {
		return err
	}

	return nil
}
