package handler

import (
	"errors"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/utils"
	"net/http"
)

const urlLength = 10

func handleBadRequest(w http.ResponseWriter, err *error) {
	w.WriteHeader(http.StatusBadRequest)

	if err != nil {
		var error = *err
		w.Write([]byte(error.Error()))
	} else {
		w.Write([]byte("bad request"))
	}
}

func handleCreatedResponse(w http.ResponseWriter, code string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("http://localhost:8080/" + code))
}

func saveUrl(db *model.Database, body []byte) (error, *string) {
	if len(body) == 0 {
		return errors.New("empty body"), nil
	}

	originalURL := string(body)
	generatedURLCode, err := utils.RandomCode(urlLength)

	if err != nil {
		return err, nil
	}

	db.Put(generatedURLCode, originalURL)
	return nil, &generatedURLCode
}

func buildSaveHandler(db *model.Database) func(w http.ResponseWriter, r *http.Request) {
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

		err, code := saveUrl(db, body)

		if err != nil {
			handleBadRequest(w, &err)
			return
		}

		handleCreatedResponse(w, *code)
	}
}
