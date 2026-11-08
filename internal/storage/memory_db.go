package storage

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"sync"
)

type MemoryDatabase struct {
	mu      sync.Mutex
	rowsMap map[string]string
}

func NewMemoryDatabase() *MemoryDatabase {
	var database = &MemoryDatabase{
		rowsMap: make(map[string]string),
	}

	return database
}

func (db *MemoryDatabase) Put(ctx context.Context, key string, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.rowsMap[key]

	if ok {
		return model.ErrDBConflict
	}

	db.rowsMap[key] = value

	return nil
}

func (db *MemoryDatabase) PutBatch(ctx context.Context, pairs []model.ValuesPair) error {
	for _, pair := range pairs {
		err := db.Put(ctx, pair.Key, pair.Value)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *MemoryDatabase) Get(ctx context.Context, key string) (*string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, ok := db.rowsMap[key]

	if ok {
		return &res, nil
	}

	return nil, model.ErrDBNotFound
}
