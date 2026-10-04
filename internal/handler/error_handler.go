package handler

import (
	"errors"
	"log"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, model.ErrDBConflict) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(http.StatusText(http.StatusConflict)))
		return
	}

	if errors.Is(err, model.ErrDBNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	if errors.Is(err, ErrInternal) {
		log.Printf("Error occured: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	log.Printf("Error occured: %v", err)

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(http.StatusText(http.StatusBadRequest)))
}
