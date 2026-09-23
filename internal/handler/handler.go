package handler

import (
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/utils"
	"net/http"
)

const urlLength = 10

func handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	originalURL := string(body)
	generatedURLCode, err := utils.RandomCode(urlLength)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	model.Put(generatedURLCode, originalURL)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("http://localhost:8080/" + generatedURLCode))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	var urlCode = r.URL.Path

	ok, res := model.Get(urlCode)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	w.WriteHeader(http.StatusPermanentRedirect)
	w.Header().Add("Location", res)
}

func HandleShortenerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleSave(w, r)
		return
	}

	if r.Method == http.MethodGet {
		handleGet(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("method is not supported"))
}
