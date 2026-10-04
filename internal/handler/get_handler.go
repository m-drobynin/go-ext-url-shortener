package handler

import (
	"fmt"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BuildGetHandler(service service.URLService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var urlCode = chi.URLParam(r, "urlCode")

		if len(urlCode) == 0 {
			handleError(w, fmt.Errorf("code is empty: %w", ErrBadRequest))
			return
		}

		res, err := service.RetrieveURL(urlCode)

		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Add("Location", *res)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
