package service

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
)

type DBProvider interface {
	Put(ctx context.Context, key string, value string) error
	PutBatch(ctx context.Context, pairs []model.ValuesPair) error
	Get(ctx context.Context, key string) (*string, error)
	GetCode(ctx context.Context, originalURL string) (*string, error)
}
