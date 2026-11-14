package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBatchShortenHandlerEmptyRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildBatchShortenHandler(&config, service)

	jsonContentType := "application/json"

	emptyRequest := []model.ShortenBatchRequestRow{}

	t.Run("empty rows", func(t *testing.T) {
		var reader io.Reader = nil

		reqJSONBytes, _ := json.Marshal(emptyRequest)
		reader = bytes.NewReader(reqJSONBytes)

		r := httptest.NewRequest(http.MethodPost, "/", reader)

		r.Header.Add("Content-Type", jsonContentType)

		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Код ответа не совпадает с ожидаемым")
	})
}

func TestBatchShortenHandlerInvalidHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildBatchShortenHandler(&config, service)

	emptyRequest := []model.ShortenBatchRequestRow{}

	t.Run("invalid header", func(t *testing.T) {
		var reader io.Reader = nil

		reqJSONBytes, _ := json.Marshal(emptyRequest)
		reader = bytes.NewReader(reqJSONBytes)

		r := httptest.NewRequest(http.MethodPost, "/", reader)
		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Код ответа не совпадает с ожидаемым")
	})
}

func TestBatchShortenHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildBatchShortenHandler(&config, service)

	jsonContentType := "application/json"
	validBodyMessageRaw, _ := json.Marshal([]model.ShortenBatchResponseRow{
		{
			CorrelationID: "id_1",
			ShortURL:      config.GetBaseURL() + "/code_1",
		},
		{
			CorrelationID: "id_2",
			ShortURL:      config.GetBaseURL() + "/code_2",
		},
	})
	validBodyMessage := string(validBodyMessageRaw)

	validRequest := []model.ShortenBatchRequestRow{
		{
			CorrelationID: "id_1",
			OriginalURL:   "url_1",
		},
		{
			CorrelationID: "id_2",
			OriginalURL:   "url_2",
		},
	}

	serviceRes := []model.ValuesPair{
		{
			Key:   "code_1",
			Value: "url_1",
		},
		{
			Key:   "code_2",
			Value: "url_2",
		},
	}

	t.Run("valid request", func(t *testing.T) {
		var reader io.Reader = nil

		reqJSONBytes, _ := json.Marshal(validRequest)
		reader = bytes.NewReader(reqJSONBytes)

		ctx := context.WithValue(t.Context(), "user_id", "user_id")
		r := httptest.NewRequestWithContext(ctx, http.MethodPost, "/", reader)

		r.Header.Add("Content-Type", jsonContentType)

		w := httptest.NewRecorder()

		service.EXPECT().SaveURLBatch(ctx, gomock.Any(), gomock.Any()).Return(&serviceRes, nil)

		handler(w, r)

		assert.Equal(t, http.StatusCreated, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.JSONEq(t, validBodyMessage, w.Body.String(), "invalid response body")
	})
}
