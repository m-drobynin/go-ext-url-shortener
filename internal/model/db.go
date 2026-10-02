package model

import (
	"sync"
)

type Database struct {
	mu      sync.Mutex
	results map[string]string
}

func NewDatabase() *Database {
	var database = &Database{}
	database.results = make(map[string]string)
	return database
}

func (db *Database) TryPut(key string, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.results[key]

	if ok {
		return ErrDBConflict
	}

	db.results[key] = value
	return nil
}

func (db *Database) Get(key string) (*string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, ok := db.results[key]

	if ok {
		return &res, nil
	}

	return nil, ErrDBNotFound
}
