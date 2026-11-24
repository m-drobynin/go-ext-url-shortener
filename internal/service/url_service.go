package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/storage"
	"m-drobynin/go-ext-url-shortener/internal/utils"

	"go.uber.org/zap"
)

const codeLength = 10
const maxSaveAttempts = 10

type URLService interface {
	SaveURL(ctx context.Context, originalURL string, userID string) (*string, error)
	SaveURLBatch(ctx context.Context, urls []string, userID string) (*[]model.ValuesPair, error)
	RetrieveURL(ctx context.Context, code string) (*model.URLRow, error)
	RetrieveCode(ctx context.Context, originalURL string) (*string, error)
	GetAllByUserID(ctx context.Context, userID string) (*[]model.URLRow, error)
	DeleteBatch(ctx context.Context, codes []string, userID string) error
}

type BatchDeleteJob struct {
	UserID string
	Codes  []string
}

type URLServiceImpl struct {
	db              storage.UrlsDBProvider
	batchDeleteJobs chan BatchDeleteJob
	log             *zap.Logger
}

func NewURLServiceImpl(db storage.UrlsDBProvider, log *zap.Logger) *URLServiceImpl {
	service := &URLServiceImpl{
		db:              db,
		batchDeleteJobs: make(chan BatchDeleteJob, 100),
		log:             log,
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

func (service *URLServiceImpl) RetrieveURL(ctx context.Context, code string) (*model.URLRow, error) {
	return service.db.Get(ctx, code)
}

func (service *URLServiceImpl) RetrieveCode(ctx context.Context, originalURL string) (*string, error) {
	return service.db.GetCode(ctx, originalURL)
}

func (service *URLServiceImpl) GetAllByUserID(ctx context.Context, userID string) (*[]model.URLRow, error) {
	return service.db.GetAllByUserID(ctx, userID)
}

func (service *URLServiceImpl) DeleteBatch(ctx context.Context, codes []string, userID string) error {
	job := BatchDeleteJob{
		UserID: userID,
		Codes:  codes,
	}

	select {
	case service.batchDeleteJobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (service *URLServiceImpl) RunBatchDeleteWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-service.batchDeleteJobs:
			if !ok {
				return
			}
			if err := service.db.DeleteBatch(ctx, job.Codes, job.UserID); err != nil {
				log.Printf("error while removing the user %v batch %v", job.UserID, job.Codes)
			}
		}
	}
}
