package model

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"uuid"
)

type FileStorage struct {
	filename string
	file     *os.File
	scanner  *bufio.Scanner
	db       *MemoryDatabase
}

func NewFileStorage(filename string) (*FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	fs := FileStorage{
		filename: filename,
		file:     file,
		scanner:  bufio.NewScanner(file),
		db:       NewMemoryDatabase(),
	}

	if err := fs.read(); err != nil {
		return nil, err
	}

	return &fs, nil
}

func (fs *FileStorage) TryPut(ctx context.Context, key string, value string) error {
	err := fs.db.TryPut(ctx, key, value)

	if err != nil {
		return err
	}

	row := ShortenedURLRow{
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

func (fs *FileStorage) Get(ctx context.Context, key string) (*string, error) {
	return fs.db.Get(ctx, key)
}

func (fs *FileStorage) read() error {
	for fs.scanner.Scan() {
		data := fs.scanner.Bytes()
		row := ShortenedURLRow{}

		if err := json.Unmarshal(data, &row); err != nil {
			return err
		}

		fs.db.rowsMap[row.ShortURL] = row.OriginalURL
	}

	return fs.scanner.Err()
}
