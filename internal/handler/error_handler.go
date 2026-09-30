package handler

import (
	"errors"
	"log"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
)

func handleError(w http.ResponseWriter, err *error) {
	log.Printf("Error occured: %v", err)

	if errors.Is(*err, model.DbConflictError) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(http.StatusText(http.StatusConflict)))
		return
	}

	if errors.Is(*err, model.DbNotFoundError) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	if errors.Is(*err, model.InternalError) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(http.StatusText(http.StatusBadRequest)))
}
