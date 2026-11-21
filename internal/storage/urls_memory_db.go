package storage

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"sync"
	"uuid"
)

type UrlsMemoryDatabase struct {
	mu            sync.Mutex
	rowsMap       map[string]model.URLRow
	reverseURLMap map[string]string
	userCodes     map[string]map[string]struct{}
}

func NewUrlsMemoryDatabase() *UrlsMemoryDatabase {
	var database = &UrlsMemoryDatabase{
		rowsMap:       make(map[string]model.URLRow),
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

	db.addRow(model.URLRow{
		UUID:        uuid.NewV4().String(),
		ShortURL:    key,
		OriginalURL: value,
		UserID:      userID,
	})

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

func (db *UrlsMemoryDatabase) Get(ctx context.Context, key string) (*model.URLRow, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rowData, ok := db.rowsMap[key]

	if ok {
		if rowData.Deleted {
			return nil, model.ErrDBDeleted
		}

		return &rowData, nil
	}

	return nil, model.ErrDBNotFound
}

func (db *UrlsMemoryDatabase) GetCode(ctx context.Context, originalURL string) (*string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	code, ok := db.reverseURLMap[originalURL]

	if ok {
		row := db.rowsMap[code]

		if row.Deleted {
			return nil, model.ErrDBDeleted
		}

		return &code, nil
	}

	return nil, model.ErrDBNotFound
}

func (db *UrlsMemoryDatabase) GetAllByUserID(ctx context.Context, userID string) (*[]model.URLRow, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	userCodes, ok := db.userCodes[userID]
	var result []model.URLRow

	if !ok {
		return &result, nil
	}

	for code := range userCodes {
		row := db.rowsMap[code]

		if !row.Deleted {
			result = append(result, row)
		}
	}

	return &result, nil
}

func (db *UrlsMemoryDatabase) DeleteBatch(ctx context.Context, codes []string, userID string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, code := range codes {
		row, ok := db.rowsMap[code]

		if ok && row.UserID == userID {
			row.Deleted = true
			db.rowsMap[code] = row
		}
	}

	return nil
}

func (db *UrlsMemoryDatabase) addRow(row model.URLRow) {
	_, ok := db.userCodes[row.UserID]

	if !ok {
		db.userCodes[row.UserID] = make(map[string]struct{})
	}

	db.rowsMap[row.ShortURL] = row
	db.reverseURLMap[row.OriginalURL] = row.ShortURL
	db.userCodes[row.UserID][row.ShortURL] = struct{}{}
}
