package handler

import (
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
)

func BuildShortenerHandler(db *model.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var saveHandler = buildSaveHandler(db)
		var getHandler = buildGetHandler(db)

		if r.Method == http.MethodPost {
			saveHandler(w, r)
			return
		}

		if r.Method == http.MethodGet {
			getHandler(w, r)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method is not supported"))
	}
}
