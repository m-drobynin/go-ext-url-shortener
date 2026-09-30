package handler

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type GetTestCase struct {
	name string
	path string

	expectedCode           int
	expectedBody           *string
	expectedLocationHeader *string
}

var existingURL = "existing_url"

type TestGetUrlServiceImpl struct {
}

func (service *TestGetUrlServiceImpl) SaveURL(originalURL string) (*string, error) {
	return nil, nil
}

func (service *TestGetUrlServiceImpl) RetrieveURL(code string) (*string, error) {
	if code == "existing_key" {
		return &existingURL, nil
	}

	return nil, model.DbNotFoundError
}

func TestGetHandler(t *testing.T) {
	service := &TestGetUrlServiceImpl{}

	handler := BuildGetHandler(service)

	var notFoundBody = http.StatusText(http.StatusNotFound)
	var badRequestBody = http.StatusText(http.StatusBadRequest)

	cases := []GetTestCase{
		{
			name: "empty url",
			path: "",

			expectedCode: http.StatusBadRequest,
			expectedBody: &badRequestBody,
		},
		{
			name: "non existing url",
			path: "nope",

			expectedCode: http.StatusNotFound,
			expectedBody: &notFoundBody,
		},
		{
			name: "existing url",
			path: "existing_key",

			expectedCode:           http.StatusTemporaryRedirect,
			expectedLocationHeader: &existingURL,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://localhost/"+tc.path, nil)
			w := httptest.NewRecorder()

			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("urlCode", tc.path)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))

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
