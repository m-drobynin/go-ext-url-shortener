package handler

import (
	"context"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var existingURL = "http://existing.com/url"
var existingURLCode = "existing_code"

type SaveTestCase struct {
	name        string
	url         *string
	contentType *string

	expectedCode int
	expectedBody *string
}

type TestSaveURLServiceImpl struct {
}

func (urlService *TestSaveURLServiceImpl) SaveURL(ctx context.Context, originalURL string) (*string, error) {
	res := "urlcodetest"

	if originalURL == existingURL {
		return nil, service.ErrUniqueValue
	}

	return &res, nil
}

func (service *TestSaveURLServiceImpl) SaveURLBatch(ctx context.Context, urls []string) (*[]model.ValuesPair, error) {
	return nil, nil
}

func (service *TestSaveURLServiceImpl) RetrieveURL(ctx context.Context, code string) (*string, error) {
	return nil, nil
}

func (service *TestSaveURLServiceImpl) RetrieveCode(ctx context.Context, originalURL string) (*string, error) {
	if originalURL == existingURL {
		return &existingURLCode, nil
	}

	return nil, nil
}

func TestSaveHandler(t *testing.T) {
	config := config.GetAppConfigDefaults()
	service := &TestSaveURLServiceImpl{}

	handler := BuildSaveHandler(&config, service)

	plainContentType := "text/plain"
	emptyStr := ""
	validURL := "http://example.com/path"
	badRequestBody := http.StatusText(http.StatusBadRequest)
	validBodyMessage := config.BaseURL

	cases := []SaveTestCase{
		{
			name: "empty url nil",

			url:         nil,
			contentType: &plainContentType,

			expectedCode: http.StatusBadRequest,
			expectedBody: &badRequestBody,
		},
		{
			name: "empty url",

			url:         &emptyStr,
			contentType: &plainContentType,

			expectedCode: http.StatusBadRequest,
			expectedBody: &badRequestBody,
		},
		{
			name: "invalid header",

			url: &validURL,

			expectedCode: http.StatusBadRequest,
			expectedBody: &badRequestBody,
		},
		{
			name:        "valid request",
			url:         &validURL,
			contentType: &plainContentType,

			expectedCode: http.StatusCreated,
			expectedBody: &validBodyMessage,
		},
		{
			name:        "conflict",
			url:         &existingURL,
			contentType: &plainContentType,

			expectedCode: http.StatusConflict,
			expectedBody: &validBodyMessage,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var reader io.Reader = nil

			if tc.url != nil {
				reader = strings.NewReader(*tc.url)
			}

			r := httptest.NewRequest(http.MethodPost, "/", reader)

			if tc.contentType != nil {
				r.Header.Add("Content-Type", *tc.contentType)
			}

			w := httptest.NewRecorder()

			handler(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Regexp(t, *tc.expectedBody, w.Body.String(), "Невалидный ответ")
		})
	}
}
