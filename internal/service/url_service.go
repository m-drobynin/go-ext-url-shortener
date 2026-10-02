package service

import (
	"m-drobynin/go-ext-url-shortener/internal/utils"
)

const urlLength = 10
const maxSaveAttempts = 10

type URLService interface {
	SaveURL(originalURL string) (*string, error)
	RetrieveURL(code string) (*string, error)
}

type URLServiceImpl struct {
	db DbProvider
}

func NewURLServiceImpl(db DbProvider) *URLServiceImpl {
	service := &URLServiceImpl{}
	service.db = db
	return service
}

func (service *URLServiceImpl) SaveURL(originalURL string) (*string, error) {
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

	return nil, ErrSaveUrlMaxAttempts
}

func (service *URLServiceImpl) RetrieveURL(code string) (*string, error) {
	return service.db.Get(code)
}
