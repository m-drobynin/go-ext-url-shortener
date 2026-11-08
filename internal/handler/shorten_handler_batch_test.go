package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BatchSaveTestCase struct {
	name         string
	rows         *[]model.ShortenBatchRequestRow
	contentType  *string
	expectedCode int
}

type TestBatchSaveURLServiceImpl struct {
}

func (service *TestBatchSaveURLServiceImpl) SaveURL(ctx context.Context, originalURL string) (*string, error) {
	return nil, nil
}

func (service *TestBatchSaveURLServiceImpl) SaveURLBatch(ctx context.Context, urls []string) (*[]model.ValuesPair, error) {
	res := []model.ValuesPair{
		{
			Key:   "code_1",
			Value: "url_1",
		},
		{
			Key:   "code_2",
			Value: "url_2",
		},
	}

	return &res, nil
}

func (service *TestBatchSaveURLServiceImpl) RetrieveURL(ctx context.Context, code string) (*string, error) {
	return nil, nil
}

func TestBatchShortenHandler(t *testing.T) {
	config := config.GetAppConfigDefaults()
	service := &TestBatchSaveURLServiceImpl{}

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

	emptyRequest := []model.ShortenBatchRequestRow{}
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

	cases := []BatchSaveTestCase{
		{
			name: "empty rows",

			rows:        &emptyRequest,
			contentType: &jsonContentType,

			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid header",
			rows: &emptyRequest,

			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "valid request",
			rows:        &validRequest,
			contentType: &jsonContentType,

			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var reader io.Reader = nil

			reqJsonBytes, _ := json.Marshal(tc.rows)
			reader = bytes.NewReader(reqJsonBytes)

			r := httptest.NewRequest(http.MethodPost, "/", reader)

			if tc.contentType != nil {
				r.Header.Add("Content-Type", *tc.contentType)
			}

			w := httptest.NewRecorder()

			handler(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")

			if len(*tc.rows) != 0 {
				assert.JSONEq(t, validBodyMessage, w.Body.String(), "invalid response body")
			}
		})
	}
}
