package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"os"
)

type UsersFileStorage struct {
	filename string
	file     *os.File
	scanner  *bufio.Scanner
	db       *UsersMemoryDatabase
}

func NewUsersFileStorage(filename string) (*UsersFileStorage, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	fs := UsersFileStorage{
		filename: filename,
		file:     file,
		scanner:  bufio.NewScanner(file),
		db:       NewUsersMemoryDatabase(),
	}

	if err := fs.read(); err != nil {
		return nil, err
	}

	return &fs, nil
}

func (fs *UsersFileStorage) Create(ctx context.Context) (*model.User, error) {
	user, err := fs.db.Create(ctx)

	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)

	if err != nil {
		return nil, err
	}

	if _, err = fs.file.Write(data); err != nil {
		return nil, err
	}

	if _, err = fs.file.WriteString("\n"); err != nil {
		return nil, err
	}

	return user, nil
}

func (fs *UsersFileStorage) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return fs.db.GetByID(ctx, userID)
}

func (fs *UsersFileStorage) read() error {
	for fs.scanner.Scan() {
		data := fs.scanner.Bytes()
		user := model.User{}

		if err := json.Unmarshal(data, &user); err != nil {
			return err
		}

		fs.db.usersMap[user.UserID] = user
	}

	return fs.scanner.Err()
}
