package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
)

func BuildDeleteHandler(service service.URLService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
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

		var codes []string

		if err = json.Unmarshal(body, &codes); err != nil {
			handleError(w, fmt.Errorf("error while unmarshalling: %v: %w", err, ErrInternal))
			return
		}

		if len(codes) == 0 {
			handleError(w, fmt.Errorf("batch is empty: %v: %w", err, ErrBadRequest))
			return
		}

		userID := r.Context().Value(model.UserIDKey).(string)
		err = service.DeleteBatch(r.Context(), codes, userID)

		if err != nil {
			handleError(w, fmt.Errorf("error while reading body: %v: %w", err, ErrInternal))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
