package service

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/utils"

	"go.uber.org/zap"
)

const urlLength = 10
const maxSaveAttempts = 10

type URLService interface {
	SaveURL(ctx context.Context, originalURL string) (*string, error)
	RetrieveURL(ctx context.Context, code string) (*string, error)
}

type URLServiceImpl struct {
	log *zap.Logger
	db  DBProvider
}

func NewURLServiceImpl(db DBProvider, log *zap.Logger) *URLServiceImpl {
	service := &URLServiceImpl{
		db:  db,
		log: log,
	}

	return service
}

func (service *URLServiceImpl) SaveURL(ctx context.Context, originalURL string) (*string, error) {
	i := 0

	for i < maxSaveAttempts {
		i++
		generatedURLCode, err := utils.RandomCode(urlLength)

		if err != nil {
			return nil, err
		}

		err = service.db.TryPut(ctx, generatedURLCode, originalURL)

		if err == nil {
			return &generatedURLCode, nil
		}
	}

	return nil, ErrSaveURLMaxAttempts
}

func (service *URLServiceImpl) RetrieveURL(ctx context.Context, code string) (*string, error) {
	return service.db.Get(ctx, code)
}
