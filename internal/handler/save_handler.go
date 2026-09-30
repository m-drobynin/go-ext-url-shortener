package handler

import (
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
)

const urlLength = 10
const maxSaveAttempts = 10

type SaveHandlerConfig interface {
	GetBaseURL() string
}

func handleCreatedResponse(config SaveHandlerConfig, w http.ResponseWriter, code string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(config.GetBaseURL() + "/" + code))
}

func BuildSaveHandler(config SaveHandlerConfig, service service.UrlService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			handleError(w, &model.BadRequestError)
			return
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			handleError(w, &err)
			return
		}

		if len(body) == 0 {
			handleError(w, &model.BadRequestError)
			return
		}

		code, err := service.SaveURL(string(body))

		if err != nil {
			handleError(w, &err)
			return
		}

		handleCreatedResponse(config, w, *code)
	}
}
