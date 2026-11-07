package model

import (
	"context"
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

func (db *MemoryDatabase) TryPut(ctx context.Context, key string, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.rowsMap[key]

	if ok {
		return ErrDBConflict
	}

	db.rowsMap[key] = value

	return nil
}

func (db *MemoryDatabase) Get(ctx context.Context, key string) (*string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, ok := db.rowsMap[key]

	if ok {
		return &res, nil
	}

	return nil, ErrDBNotFound
}
