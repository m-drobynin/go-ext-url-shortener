package handler

import (
	"context"
	"encoding/json"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUrlsHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockURLService(ctrl)

	handler := BuildUserUrlsHandler(service)

	t.Run("get user urls", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), model.UserIDKey, "user_id"))

		service.EXPECT().GetAllByUserID(gomock.Any(), "user_id").Return(&[]model.ValuesPair{
			{
				Key:   "url_code_1",
				Value: "original_url_1",
			},
		}, nil)

		handler(w, r)

		assert.Equal(t, http.StatusOK, w.Code, "Код ответа не совпадает с ожидаемым")

		var result []model.UrlsPairRow

		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, result[0].ShortURL, "url_code_1")
		assert.Equal(t, result[0].OriginalURL, "original_url_1")
	})
}
