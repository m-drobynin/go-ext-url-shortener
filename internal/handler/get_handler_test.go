package handler

import (
	"context"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupTestGetCall(path string) (*httptest.ResponseRecorder, *http.Request) {
	r := httptest.NewRequest(http.MethodGet, "http://localhost/"+path, nil)
	w := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("urlCode", path)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))

	return w, r
}

func TestGetHandlerEmptyUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockURLService(ctrl)

	handler := BuildGetHandler(service)

	t.Run("empty url", func(t *testing.T) {
		w, r := setupTestGetCall("")

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Equal(t, http.StatusText(http.StatusBadRequest), w.Body.String(), "Тело ответа не совпадает с ожидаемым")
	})
}

func TestGetHandlerNonExistingUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockURLService(ctrl)

	handler := BuildGetHandler(service)

	t.Run("non existing url", func(t *testing.T) {
		service.EXPECT().RetrieveURL(gomock.Any(), gomock.Any()).Return(nil, model.ErrDBNotFound)

		w, r := setupTestGetCall("nope")
		handler(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code, "Код ответа не совпадает с ожидаемым")
		assert.Equal(t, http.StatusText(http.StatusNotFound), w.Body.String(), "Тело ответа не совпадает с ожидаемым")
	})
}

func TestGetHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockURLService(ctrl)

	handler := BuildGetHandler(service)

	t.Run("existing url", func(t *testing.T) {
		service.EXPECT().RetrieveURL(gomock.Any(), gomock.Any()).Return(&existingURL, nil)

		w, r := setupTestGetCall("existing_key")
		handler(w, r)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code, "Код ответа не совпадает с ожидаемым")
		var locationHeader = w.Result().Header.Get("Location")
		assert.Equal(t, existingURL, locationHeader, "Заголовок Location не совпадает с ожидаемым")
	})
}
