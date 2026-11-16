package handler

import (
	"encoding/json"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
)

func BuildUserUrlsHandler(config SaveHandlerConfig, urlService service.URLService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(model.UserIDKey).(string)

		userRows, err := urlService.GetAllByUserID(r.Context(), userID)

		if err != nil {
			handleError(w, err)
			return
		}

		if len(*userRows) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var rows []model.UrlsPairRow

		for _, pair := range *userRows {
			rows = append(rows, model.UrlsPairRow{
				ShortURL:    config.GetBaseURL() + "/" + pair.Key,
				OriginalURL: pair.Value,
			})
		}

		res, err := json.Marshal(rows)

		if err != nil {
			handleError(w, err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}
