package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"os"
	"uuid"
)

type UrlsFileStorage struct {
	filename string
	file     *os.File
	scanner  *bufio.Scanner
	db       *UrlsMemoryDatabase
}

func NewUrlsFileStorage(filename string) (*UrlsFileStorage, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	fs := UrlsFileStorage{
		filename: filename,
		file:     file,
		scanner:  bufio.NewScanner(file),
		db:       NewUrlsMemoryDatabase(),
	}

	if err := fs.read(); err != nil {
		return nil, err
	}

	return &fs, nil
}

func (fs *UrlsFileStorage) Put(ctx context.Context, key string, value string) error {
	err := fs.db.Put(ctx, key, value)

	if err != nil {
		return err
	}

	row := model.ShortenedURLRow{
		ShortURL:    key,
		OriginalURL: value,
		UUID:        uuid.NewV4().String(),
	}

	data, err := json.Marshal(row)

	if err != nil {
		return err
	}

	if _, err = fs.file.Write(data); err != nil {
		return err
	}

	if _, err = fs.file.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

func (fs *UrlsFileStorage) PutBatch(ctx context.Context, pairs []model.ValuesPair) error {
	for _, pair := range pairs {
		err := fs.Put(ctx, pair.Key, pair.Value)

		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *UrlsFileStorage) Get(ctx context.Context, key string) (*string, error) {
	return fs.db.Get(ctx, key)
}

func (fs *UrlsFileStorage) GetCode(ctx context.Context, originalURL string) (*string, error) {
	return fs.db.GetCode(ctx, originalURL)
}

func (fs *UrlsFileStorage) read() error {
	for fs.scanner.Scan() {
		data := fs.scanner.Bytes()
		row := model.ShortenedURLRow{}

		if err := json.Unmarshal(data, &row); err != nil {
			return err
		}

		fs.db.rowsMap[row.ShortURL] = row.OriginalURL
	}

	return fs.scanner.Err()
}
