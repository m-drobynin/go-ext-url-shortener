package storage

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"sync"
)

type UrlsMemoryDatabase struct {
	mu            sync.Mutex
	rowsMap       map[string]string
	reverseURLMap map[string]string
	userCodes     map[string]map[string]struct{}
}

func NewUrlsMemoryDatabase() *UrlsMemoryDatabase {
	var database = &UrlsMemoryDatabase{
		rowsMap:       make(map[string]string),
		reverseURLMap: make(map[string]string),
		userCodes:     make(map[string]map[string]struct{}),
	}

	return database
}

func (db *UrlsMemoryDatabase) Put(ctx context.Context, key string, value string, userID string) error {
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

	_, ok = db.userCodes[userID]

	if !ok {
		db.userCodes[userID] = make(map[string]struct{})
	}

	db.rowsMap[key] = value
	db.reverseURLMap[value] = key
	db.userCodes[userID][key] = struct{}{}

	return nil
}

func (db *UrlsMemoryDatabase) PutBatch(ctx context.Context, pairs []model.ValuesPair, userID string) error {
	for _, pair := range pairs {
		err := db.Put(ctx, pair.Key, pair.Value, userID)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *UrlsMemoryDatabase) Get(ctx context.Context, key string) (*string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, ok := db.rowsMap[key]

	if ok {
		return &res, nil
	}

	return nil, model.ErrDBNotFound
}

func (db *UrlsMemoryDatabase) GetCode(ctx context.Context, originalURL string) (*string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, ok := db.reverseURLMap[originalURL]

	if ok {
		return &res, nil
	}

	return nil, model.ErrDBNotFound
}

func (db *UrlsMemoryDatabase) GetAllByUserId(ctx context.Context, userID string) (*[]model.ValuesPair, error) {
	userCodes, ok := db.userCodes[userID]
	var result []model.ValuesPair

	if !ok {
		return &result, nil
	}

	for code := range userCodes {
		result = append(result, model.ValuesPair{
			Key:   code,
			Value: db.rowsMap[code],
		})
	}

	return &result, nil
}
