package storage

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"sync"
)

type MemoryDatabase struct {
	mu            sync.Mutex
	rowsMap       map[string]string
	reverseURLMap map[string]string
}

func NewMemoryDatabase() *MemoryDatabase {
	var database = &MemoryDatabase{
		rowsMap:       make(map[string]string),
		reverseURLMap: make(map[string]string),
	}

	return database
}

func (db *MemoryDatabase) Put(ctx context.Context, key string, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.reverseURLMap[value]

	if ok {
		return model.ErrUrlsConflict
	}

	_, ok = db.rowsMap[key]

	if ok {
		return model.ErrCodeConflict
	}

	db.rowsMap[key] = value
	db.reverseURLMap[value] = key

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

func (db *MemoryDatabase) GetCode(ctx context.Context, originalURL string) (*string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, ok := db.reverseURLMap[originalURL]

	if ok {
		return &res, nil
	}

	return nil, model.ErrDBNotFound
}
