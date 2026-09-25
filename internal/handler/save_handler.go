package handler

import (
	"errors"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/utils"
	"net/http"
)

const urlLength = 10

func handleCreatedResponse(w http.ResponseWriter, code string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("http://localhost:8080/" + code))
}

func saveURL(db *model.Database, body []byte) (*string, error) {
	if len(body) == 0 {
		return nil, errors.New("empty body")
	}

	originalURL := string(body)
	generatedURLCode, err := utils.RandomCode(urlLength)

	if err != nil {
		return nil, err
	}

	db.Put(generatedURLCode, originalURL)
	return &generatedURLCode, nil
}

func BuildSaveHandler(db *model.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			handleBadRequest(w, nil)
			return
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			handleBadRequest(w, &err)
			return
		}

		code, err := saveURL(db, body)

		if err != nil {
			handleBadRequest(w, &err)
			return
		}

		handleCreatedResponse(w, *code)
	}
}
