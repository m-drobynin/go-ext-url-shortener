package service

import "context"

type DBProvider interface {
	TryPut(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (*string, error)
}
