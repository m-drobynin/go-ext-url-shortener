package handler

import (
	"context"
	"encoding/json"
	"m-drobynin/go-ext-url-shortener/internal/config"
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

	config := config.GetAppConfigDefaults()
	service := mocks.NewMockURLService(ctrl)

	handler := BuildUserUrlsHandler(&config, service)

	t.Run("get user urls", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, config.GetBaseURL(), nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), model.UserIDKey, "user_id"))

		service.EXPECT().GetAllByUserID(gomock.Any(), "user_id").Return(&[]model.URLRow{
			{
				ShortURL:    "url_code_1",
				OriginalURL: "original_url_1",
			},
		}, nil)

		handler(w, r)

		assert.Equal(t, http.StatusOK, w.Code, "Код ответа не совпадает с ожидаемым")

		var result []model.UrlsPairRow

		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, result[0].ShortURL, config.GetBaseURL()+"/url_code_1")
		assert.Equal(t, result[0].OriginalURL, "original_url_1")
	})
}
