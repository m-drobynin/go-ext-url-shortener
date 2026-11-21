package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"os"
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

func (fs *UrlsFileStorage) Put(ctx context.Context, key string, value string, userID string) error {
	err := fs.db.Put(ctx, key, value, userID)

	if err != nil {
		return err
	}

	row, err := fs.db.Get(ctx, key)

	if err != nil {
		return err
	}

	return fs.writeRow(row)
}

func (fs *UrlsFileStorage) PutBatch(ctx context.Context, pairs []model.ValuesPair, userID string) error {
	for _, pair := range pairs {
		err := fs.Put(ctx, pair.Key, pair.Value, userID)

		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *UrlsFileStorage) Get(ctx context.Context, key string) (*model.URLRow, error) {
	return fs.db.Get(ctx, key)
}

func (fs *UrlsFileStorage) GetCode(ctx context.Context, originalURL string) (*string, error) {
	return fs.db.GetCode(ctx, originalURL)
}

func (fs *UrlsFileStorage) GetAllByUserID(ctx context.Context, userID string) (*[]model.URLRow, error) {
	return fs.db.GetAllByUserID(ctx, userID)
}

func (fs *UrlsFileStorage) DeleteBatch(ctx context.Context, codes []string, userID string) error {
	for _, code := range codes {
		row, err := fs.db.Get(ctx, code)

		if err != nil {
			return err
		}

		if row.UserID == userID {
			row.Deleted = true

			err = fs.writeRow(row)

			if err != nil {
				return err
			}
		}
	}

	return fs.db.DeleteBatch(ctx, codes, userID)
}

func (fs *UrlsFileStorage) read() error {
	for fs.scanner.Scan() {
		data := fs.scanner.Bytes()
		row := model.URLRow{}

		if err := json.Unmarshal(data, &row); err != nil {
			return err
		}

		fs.db.addRow(row)
	}

	return fs.scanner.Err()
}

func (fs *UrlsFileStorage) writeRow(row *model.URLRow) error {
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
