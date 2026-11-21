package handler

import (
	"bytes"
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

func TestDeleteUrlsHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := config.GetAppConfigDefaults()
	service := mocks.NewMockURLService(ctrl)

	handler := BuildDeleteHandler(service)
	codes := []string{"url_code"}
	userID := "user_id"

	t.Run("delete urls", func(t *testing.T) {
		reqJSONBytes, _ := json.Marshal(codes)
		reader := bytes.NewReader(reqJSONBytes)

		ctx := context.WithValue(t.Context(), model.UserIDKey, userID)
		r := httptest.NewRequestWithContext(ctx, http.MethodDelete, config.GetBaseURL(), reader)
		r.Header.Add("Content-Type", jsonContentType)
		w := httptest.NewRecorder()

		service.EXPECT().DeleteBatch(gomock.Any(), codes, userID).Return(nil)

		handler(w, r)

		assert.Equal(t, http.StatusAccepted, w.Code, "Код ответа не совпадает с ожидаемым")
	})
}
