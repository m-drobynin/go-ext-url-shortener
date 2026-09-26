package handler

import (
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type GetTestCase struct {
	name string
	url  string

	expectedCode           int
	expectedBody           *string
	expectedLocationHeader *string
}

func TestGetHandler(t *testing.T) {
	database := model.CreateDatabase()
	existingURL := "existing_url"
	database.Put("existing_key", existingURL)

	config := config.GetAppConfigDefaults()
	handler := BuildGetHandler(&config, &database)

	var notFoundBody = "not found"

	cases := []GetTestCase{
		{
			name: "empty url",
			url:  config.BaseURL,

			expectedCode: http.StatusNotFound,
			expectedBody: &notFoundBody,
		},
		{
			name: "non existing url",
			url:  config.BaseURL + "/nope",

			expectedCode: http.StatusNotFound,
			expectedBody: &notFoundBody,
		},
		{
			name: "existing url",
			url:  config.BaseURL + "/existing_key",

			expectedCode:           http.StatusTemporaryRedirect,
			expectedLocationHeader: &existingURL,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()

			handler(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")

			if tc.expectedBody != nil {
				assert.Equal(t, *tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}

			if tc.expectedLocationHeader != nil {
				var locationHeader = w.Result().Header.Get("Location")

				assert.Equal(t, *tc.expectedLocationHeader, locationHeader, "Заголовок Location не совпадает с ожидаемым")
			}
		})
	}
}
