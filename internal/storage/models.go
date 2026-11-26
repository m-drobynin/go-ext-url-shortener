package storage

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
)

type UrlsDBProvider interface {
	Put(ctx context.Context, key string, value string, userID string) error
	PutBatch(ctx context.Context, pairs []model.ValuesPair, userID string) error
	Get(ctx context.Context, key string) (*model.URLRow, error)
	GetCode(ctx context.Context, originalURL string) (*string, error)
	GetAllByUserID(ctx context.Context, userID string) (*[]model.URLRow, error)
	DeleteBatch(ctx context.Context, codes []string, userID string) error
}

type UsersDBProvider interface {
	Create(ctx context.Context) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
}
