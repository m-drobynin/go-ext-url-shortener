package service

import (
	"m-drobynin/go-ext-url-shortener/internal/logger"
	"m-drobynin/go-ext-url-shortener/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestDeleteBatchWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := mocks.NewMockUrlsDBProvider(ctrl)

	log, _ := logger.NewLogger("info")
	service := NewURLServiceImpl(storage, log)

	user1 := "user_1"
	user2 := "user_2"
	codes1 := []string{"code_11", "code_12"}
	codes2 := []string{"code_21", "code_22"}

	storage.EXPECT().DeleteBatch(gomock.Any(), codes1, user1)
	storage.EXPECT().DeleteBatch(gomock.Any(), codes2, user2)

	service.DeleteBatch(t.Context(), codes1, user1)
	service.DeleteBatch(t.Context(), codes2, user2)

	close(service.batchDeleteJobs)
	service.RunBatchDeleteWorker(t.Context())
}
