package handler

import (
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
)

func BuildGetHandler(db *model.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var urlCode = r.URL.Path

		if len(r.URL.Path) < 2 {
			handleNotFound(w)
			return
		}

		ok, res := db.Get(urlCode[1:])

		if !ok {
			handleNotFound(w)
			return
		}

		w.Header().Add("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
