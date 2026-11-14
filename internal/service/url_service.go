package service

import (
	"context"
	"errors"
	"fmt"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/storage"
	"m-drobynin/go-ext-url-shortener/internal/utils"
)

const codeLength = 10
const maxSaveAttempts = 10

type URLService interface {
	SaveURL(ctx context.Context, originalURL string, userID string) (*string, error)
	SaveURLBatch(ctx context.Context, urls []string, userID string) (*[]model.ValuesPair, error)
	RetrieveURL(ctx context.Context, code string) (*string, error)
	RetrieveCode(ctx context.Context, originalURL string) (*string, error)
}

type URLServiceImpl struct {
	db storage.UrlsDBProvider
}

func NewURLServiceImpl(db storage.UrlsDBProvider) *URLServiceImpl {
	service := &URLServiceImpl{
		db: db,
	}

	return service
}

func (service *URLServiceImpl) SaveURL(ctx context.Context, originalURL string, userID string) (*string, error) {
	for i := 0; i < maxSaveAttempts; i++ {
		generatedURLCode, err := utils.RandomCode(codeLength)

		if err != nil {
			return nil, err
		}

		err = service.db.Put(ctx, generatedURLCode, originalURL, userID)

		if err == nil {
			return &generatedURLCode, nil
		}

		if errors.Is(err, model.ErrUrlsConflict) {
			return nil, fmt.Errorf("%w:%w", ErrUniqueValue, err)
		}
	}

	return nil, ErrSaveURLMaxAttempts
}

func (service *URLServiceImpl) SaveURLBatch(ctx context.Context, urls []string, userID string) (*[]model.ValuesPair, error) {
	var results []model.ValuesPair

	for _, url := range urls {
		code, err := utils.RandomCode(codeLength)

		if err != nil {
			return nil, err
		}

		results = append(results, model.ValuesPair{
			Key:   code,
			Value: url,
		})
	}

	err := service.db.PutBatch(ctx, results, userID)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

func (service *URLServiceImpl) RetrieveURL(ctx context.Context, code string) (*string, error) {
	return service.db.Get(ctx, code)
}

func (service *URLServiceImpl) RetrieveCode(ctx context.Context, originalURL string) (*string, error) {
	return service.db.GetCode(ctx, originalURL)
}
