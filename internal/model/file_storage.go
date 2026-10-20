package model

import (
	"encoding/json"
	"os"
	"uuid"
)

type FileStorage struct {
	filename string
	rows     []ShortenedUrlRow
}

type ShortenedUrlRow struct {
	UUID        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

func NewFileStorage(filename string) (*FileStorage, error) {
	fs := FileStorage{
		filename: filename,
		rows:     []ShortenedUrlRow{},
	}

	if err := fs.read(); err != nil {
		return nil, err
	}

	return &fs, nil
}

func (fs *FileStorage) TryPut(key string, value string) error {
	newRow := ShortenedUrlRow{
		ShortUrl:    key,
		OriginalUrl: value,
		UUID:        uuid.NewV4().String(),
	}

	fs.rows = append(fs.rows, newRow)

	return fs.save()
}

func (fs *FileStorage) Get(key string) (*string, error) {
	for _, v := range fs.rows {
		if v.ShortUrl == key {
			return &v.OriginalUrl, nil
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
