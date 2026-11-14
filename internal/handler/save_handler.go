package handler

import (
	"errors"
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

func handleConflictResponse(config SaveHandlerConfig, w http.ResponseWriter, code string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusConflict)
	w.Write([]byte(config.GetBaseURL() + "/" + code))
}

func BuildSaveHandler(config SaveHandlerConfig, urlService service.URLService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			handleError(w, ErrBadRequest)
			return
		}

		bodyReader := http.MaxBytesReader(w, r.Body, maxBodyBytes)
		body, err := io.ReadAll(bodyReader)

		if err != nil {
			handleError(w, fmt.Errorf("error while reading body: %v: %w", err, ErrInternal))
			return
		}

		if len(body) == 0 {
			handleError(w, fmt.Errorf("body is empty: %v: %w", err, ErrBadRequest))
			return
		}

		originalURL := string(body)
		userId := r.Context().Value("user_id").(string)

		code, err := urlService.SaveURL(r.Context(), originalURL, userId)

		if err != nil {
			if errors.Is(err, service.ErrUniqueValue) {
				code, err = urlService.RetrieveCode(r.Context(), originalURL)

				if err != nil {
					handleError(w, fmt.Errorf("error while searching for url code: %v: %w", err, ErrInternal))
					return
				}

				handleConflictResponse(config, w, *code)
				return
			}

			handleError(w, fmt.Errorf("error while saving url: %v: %w", err, ErrInternal))
			return
		}

		handleCreatedResponse(config, w, *code)
	}
}
