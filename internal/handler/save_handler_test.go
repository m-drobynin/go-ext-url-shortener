package handler

import (
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"m-drobynin/go-ext-url-shortener/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
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

func TestSaveHandlerEmptyUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildSaveHandler(&config, service)

	plainContentType := "text/plain"

	t.Run("empty url", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		r.Header.Add("Content-Type", plainContentType)

		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Regexp(t, http.StatusText(http.StatusBadRequest), w.Body.String(), "Невалидный ответ")
	})
}

func TestSaveHandlerInvalidHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildSaveHandler(&config, service)

	t.Run("invalid header", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com/path"))
		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Regexp(t, http.StatusText(http.StatusBadRequest), w.Body.String(), "Невалидный ответ")
	})
}

func TestSaveHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildSaveHandler(&config, service)

	originalURL := "http://example.com/path"
	resCode := "example_code"

	t.Run("valid request", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
		r.Header.Add("Content-Type", "text/plain")

		service.EXPECT().SaveURL(gomock.Any(), originalURL).Return(&resCode, nil)

		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusCreated, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Equal(t, config.GetBaseURL()+"/"+resCode, w.Body.String(), "Невалидный ответ")
	})
}

func TestSaveHandlerConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildSaveHandler(&config, mockService)

	t.Run("conflict", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(existingURL))
		r.Header.Add("Content-Type", "text/plain")

		mockService.EXPECT().SaveURL(gomock.Any(), existingURL).Return(nil, service.ErrUniqueValue)
		mockService.EXPECT().RetrieveCode(gomock.Any(), existingURL).Return(&existingURLCode, nil)

		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusConflict, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Equal(t, config.GetBaseURL()+"/"+existingURLCode, w.Body.String(), "Невалидный ответ")
	})
}
