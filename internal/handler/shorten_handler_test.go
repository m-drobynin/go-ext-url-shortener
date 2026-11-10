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

func prepareShortenRequest(originalURL string) *http.Request {
	reqBody := `{"url":"` + originalURL + `"}`
	return httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
}

func prepareShortenResponse(config config.AppConfig, code string) string {
	return `{"result":"` + config.GetBaseURL() + `/` + code + `"}`
}

func TestShortenHandlerEmptyUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildShortenHandler(&config, service)

	plainContentType := "text/plain"

	t.Run("empty url", func(t *testing.T) {
		r := prepareShortenRequest("")

		r.Header.Add("Content-Type", plainContentType)

		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Regexp(t, http.StatusText(http.StatusBadRequest), w.Body.String(), "Невалидный ответ")
	})
}

func TestShortenHandlerInvalidHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildShortenHandler(&config, service)

	t.Run("invalid header", func(t *testing.T) {
		r := prepareShortenRequest("http://example.com/path")
		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Regexp(t, http.StatusText(http.StatusBadRequest), w.Body.String(), "Невалидный ответ")
	})
}

func TestShortenHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildShortenHandler(&config, service)

	originalURL := "http://example.com/path"
	resCode := "example_code"

	t.Run("valid request", func(t *testing.T) {
		r := prepareShortenRequest(originalURL)
		r.Header.Add("Content-Type", "application/json")

		service.EXPECT().SaveURL(gomock.Any(), originalURL).Return(&resCode, nil)

		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusCreated, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Equal(t, prepareShortenResponse(config, resCode), w.Body.String(), "Невалидный ответ")
	})
}

func TestShortenHandlerConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockURLService(ctrl)

	config := config.GetAppConfigDefaults()

	handler := BuildShortenHandler(&config, mockService)

	t.Run("conflict", func(t *testing.T) {
		r := prepareShortenRequest(existingURL)
		r.Header.Add("Content-Type", "application/json")

		mockService.EXPECT().SaveURL(gomock.Any(), existingURL).Return(nil, service.ErrUniqueValue)
		mockService.EXPECT().RetrieveCode(gomock.Any(), existingURL).Return(&existingURLCode, nil)

		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equal(t, http.StatusConflict, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Equal(t, prepareShortenResponse(config, existingURLCode), w.Body.String(), "Невалидный ответ")
	})
}
