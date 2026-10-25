package model

import (
	"bufio"
	"encoding/json"
	"os"
)

type FileStorage struct {
	filename string
	file     *os.File
	scanner  *bufio.Scanner
	db       *Database
}

func NewFileStorage(filename string) (*FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	fs := FileStorage{
		filename: filename,
		file:     file,
		scanner:  bufio.NewScanner(file),
		db:       NewDatabase(),
	}

	if err := fs.read(); err != nil {
		return nil, err
	}

	return &fs, nil
}

func (fs *FileStorage) TryPut(key string, value string) error {
	err := fs.db.TryPut(key, value)

	if err != nil {
		return err
	}

	row := fs.db.rows[len(fs.db.rows)-1]

	data, err := json.Marshal(row)

	if err != nil {
		return err
	}

	return os.WriteFile(fs.filename, data, 0666)
}

func (fs *FileStorage) Get(key string) (*string, error) {
	return fs.db.Get(key)
}

func (fs *FileStorage) read() error {
	for res := fs.scanner.Scan(); res; {
		data := fs.scanner.Bytes()
		row := ShortenedURLRow{}

		if err := json.Unmarshal(data, &row); err != nil {
			return err
		}

		fs.db.rows = append(fs.db.rows, row)
	}

	return nil
}
