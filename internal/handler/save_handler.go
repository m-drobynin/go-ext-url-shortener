package handler

import (
	"fmt"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
)

type SaveHandlerConfig interface {
	GetBaseURL() string
}

func handleCreatedResponse(config SaveHandlerConfig, w http.ResponseWriter, code string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(config.GetBaseURL() + "/" + code))
}

func BuildSaveHandler(config SaveHandlerConfig, service service.URLService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			handleError(w, ErrBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			handleError(w, fmt.Errorf("error while reading body: %v: %w", err, ErrInternal))
			return
		}

		if len(body) == 0 {
			handleError(w, fmt.Errorf("body is empty: %v: %w", err, ErrBadRequest))
			return
		}

		code, err := service.SaveURL(string(body))

		if err != nil {
			handleError(w, fmt.Errorf("error while saving url: %v: %w", err, ErrInternal))
			return
		}

		handleCreatedResponse(config, w, *code)
	}
}
