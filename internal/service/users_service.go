package service

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/storage"
)

type UsersService interface {
	Create(ctx context.Context) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
}

type UsersServiceImpl struct {
	db storage.UsersDBProvider
}

func NewUsersServiceImpl(db storage.UsersDBProvider) *UsersServiceImpl {
	service := &UsersServiceImpl{
		db: db,
	}

	return service
}

func (db *UsersServiceImpl) Create(ctx context.Context) (*model.User, error) {
	return db.Create(ctx)
}

func (db *UsersServiceImpl) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return db.GetByID(ctx, userID)
}
