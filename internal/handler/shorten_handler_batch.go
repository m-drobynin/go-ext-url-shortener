package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
)

func handleBatchCreatedJSONResponse(w http.ResponseWriter, rows []model.ShortenBatchResponseRow) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	res, err := json.Marshal(rows)

	if err != nil {
		return err
	}

	w.Write(res)
	return nil
}

func saveBatch(ctx context.Context, config SaveHandlerConfig, service service.URLService, request []model.ShortenBatchRequestRow) (*[]model.ShortenBatchResponseRow, error) {
	var requestMap = make(map[string]string)
	var urls = make([]string, len(request))

	for _, pair := range request {
		requestMap[pair.OriginalURL] = pair.CorrelationID
		urls = append(urls, pair.OriginalURL)
	}

	pairs, err := service.SaveURLBatch(ctx, urls)

	if err != nil {
		return nil, err
	}

	var response []model.ShortenBatchResponseRow

	for _, pair := range *pairs {
		response = append(response, model.ShortenBatchResponseRow{
			CorrelationID: requestMap[pair.Value],
			ShortURL:      config.GetBaseURL() + "/" + pair.Key,
		})
	}

	return &response, nil
}

func BuildBatchShortenHandler(config SaveHandlerConfig, service service.URLService) func(w http.ResponseWriter, r *http.Request) {
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

		var request []model.ShortenBatchRequestRow

		if err = json.Unmarshal(body, &request); err != nil {
			handleError(w, fmt.Errorf("error while unmarshalling: %v: %w", err, ErrInternal))
			return
		}

		if len(request) == 0 {
			handleError(w, fmt.Errorf("batch is empty: %v: %w", err, ErrBadRequest))
			return
		}

		responseRows, err := saveBatch(r.Context(), config, service, request)

		if err != nil {
			handleError(w, fmt.Errorf("error while batch saving: %v: %w", err, ErrInternal))
			return
		}

		err = handleBatchCreatedJSONResponse(w, *responseRows)

		if err != nil {
			handleError(w, fmt.Errorf("error while response marshalling: %v: %w", err, ErrInternal))
			return
		}
	}
}
