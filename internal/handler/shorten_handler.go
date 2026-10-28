package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
)

func handleCreatedJSONResponse(config SaveHandlerConfig, w http.ResponseWriter, code string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var response = model.ShortenResponse{
		Result: config.GetBaseURL() + "/" + code,
	}

	res, err := json.Marshal(response)

	if err != nil {
		return err
	}

	w.Write(res)
	return nil
}

func BuildShortenHandler(config SaveHandlerConfig, service service.URLService) func(w http.ResponseWriter, r *http.Request) {
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

		var request model.ShortenRequest

		if err = json.Unmarshal(body, &request); err != nil {
			handleError(w, fmt.Errorf("error while unmarshalling: %v: %w", err, ErrInternal))
			return
		}

		if len(request.URL) == 0 {
			handleError(w, fmt.Errorf("url is empty: %v: %w", err, ErrBadRequest))
			return
		}

		code, err := service.SaveURL(request.URL)

		if err != nil {
			handleError(w, fmt.Errorf("error while saving url: %v: %w", err, ErrInternal))
			return
		}

		err = handleCreatedJSONResponse(config, w, *code)

		if err != nil {
			handleError(w, fmt.Errorf("error while response marshalling: %v: %w", err, ErrInternal))
			return
		}
	}
}
