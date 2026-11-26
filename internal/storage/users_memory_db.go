package storage

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"sync"
	"uuid"
)

type UsersMemoryDatabase struct {
	mu       sync.Mutex
	usersMap map[string]model.User
}

func NewUsersMemoryDatabase() *UsersMemoryDatabase {
	var database = &UsersMemoryDatabase{
		usersMap: make(map[string]model.User),
	}

	return database
}

func (db *UsersMemoryDatabase) Create(ctx context.Context) (*model.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	user := model.User{
		UserID: uuid.NewV4().String(),
	}

	db.usersMap[user.UserID] = user

	return &user, nil
}

func (db *UsersMemoryDatabase) GetByID(ctx context.Context, userID string) (*model.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, ok := db.usersMap[userID]

	if ok {
		return &res, nil
	}

	return nil, model.ErrDBNotFound
}
