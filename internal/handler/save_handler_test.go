package handler

import (
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
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

func TestSaveHandler(t *testing.T) {
	var database = model.CreateDatabase()
	var existingUrl = "existing_url"
	database.Put("existing_key", existingUrl)

	var handler = buildSaveHandler(&database)

	var plainContentType = "text/plain"
	var emptyStr = ""
	var validUrl = "http://example.com/path"
	var emptyBodyMessage = "empty body"
	var badRequestBodyMessage = "bad request"
	var validBodyMessage = "localhost:8080"

	cases := []SaveTestCase{
		{
			name: "empty url nil",

			url:         nil,
			contentType: &plainContentType,

			expectedCode: http.StatusBadRequest,
			expectedBody: &emptyBodyMessage,
		},
		{
			name: "empty url",

			url:         &emptyStr,
			contentType: &plainContentType,

			expectedCode: http.StatusBadRequest,
			expectedBody: &emptyBodyMessage,
		},
		{
			name: "invalid header",

			url: &validUrl,

			expectedCode: http.StatusBadRequest,
			expectedBody: &badRequestBodyMessage,
		},
		{
			name:        "valid request",
			url:         &validUrl,
			contentType: &plainContentType,

			expectedCode: http.StatusCreated,
			expectedBody: &validBodyMessage,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// var bodyReader = tc.url == nil ? nil : strings.NewReader(tc.url)
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

			assert.True(t, strings.Contains(w.Body.String(), *tc.expectedBody), "Невалидный ответ")
		})
	}
}
