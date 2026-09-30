package handler

import (
	"io"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SaveTestCase struct {
	name        string
	url         *string
	contentType *string

	expectedCode int
	expectedBody *string
}

type TestSaveUrlServiceImpl struct {
}

func (service *TestSaveUrlServiceImpl) SaveURL(originalURL string) (*string, error) {
	res := "save url temp"
	return &res, nil
}

func (service *TestSaveUrlServiceImpl) RetrieveURL(code string) (*string, error) {
	return nil, nil
}

func TestSaveHandler(t *testing.T) {
	config := config.GetAppConfigDefaults()
	service := &TestSaveUrlServiceImpl{}

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
