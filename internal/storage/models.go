package storage

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
)

type UrlsDBProvider interface {
	Put(ctx context.Context, key string, value string) error
	PutBatch(ctx context.Context, pairs []model.ValuesPair) error
	Get(ctx context.Context, key string) (*string, error)
	GetCode(ctx context.Context, originalURL string) (*string, error)
}

type UsersDBProvider interface {
	Create(ctx context.Context) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
}
