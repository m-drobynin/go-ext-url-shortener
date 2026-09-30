package service

import (
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/utils"
)

const urlLength = 10
const maxSaveAttempts = 10

type UrlService interface {
	SaveURL(originalURL string) (*string, error)
	RetrieveURL(code string) (*string, error)
}

type UrlServiceImpl struct {
	db *model.Database
}

func NewUrlServiceImpl(db *model.Database) *UrlServiceImpl {
	service := &UrlServiceImpl{}
	service.db = db
	return service
}

func (service *UrlServiceImpl) SaveURL(originalURL string) (*string, error) {
	i := 0

	for i < maxSaveAttempts {
		i++
		generatedURLCode, err := utils.RandomCode(urlLength)

		if err != nil {
			return nil, err
		}

		err = service.db.TryPut(generatedURLCode, originalURL)

		if err == nil {
			return &generatedURLCode, nil
		}
	}

	return nil, model.InternalError
}

func (service *UrlServiceImpl) RetrieveURL(code string) (*string, error) {
	return service.db.Get(code)
}
