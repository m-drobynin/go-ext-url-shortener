package model

import (
	"bufio"
	"encoding/json"
	"os"
	"uuid"
)

type FileStorage struct {
	filename string
	rows     []ShortenedURLRow
	file     *os.File
	scanner  *bufio.Scanner
}

type ShortenedURLRow struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(filename string) (*FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	fs := FileStorage{
		filename: filename,
		file:     file,
		rows:     []ShortenedURLRow{},
		scanner:  bufio.NewScanner(file),
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
	if !fs.scanner.Scan() {
		return fs.scanner.Err()
	}

	data := fs.scanner.Bytes()

	if err := json.Unmarshal(data, &fs.rows); err != nil {
		return err
	}

	return nil
}
