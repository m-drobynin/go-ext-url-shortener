package middleware

import (
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

var mockHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func TestJwtAuthorizationMiddlewareSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockUsersService(ctrl)

	r := chi.NewRouter()
	r.Get("/", NewJwtHandler(mockHandler, service))

	ts := httptest.NewServer(r)
	defer ts.Close()

	userID := "test_user_id"
	jwtToken, err := buildJWT(userID)
	assert.NoError(t, err, "error while building jwt")

	req := resty.New().R()
	req.Method = "GET"
	req.URL = ts.URL
	req.Cookies = []*http.Cookie{
		{
			Name:  authCookieName,
			Value: jwtToken,
		},
	}

	service.EXPECT().GetByID(gomock.Any(), userID).Return(&model.User{
		UserID: userID,
	}, nil)

	resp, err := req.Send()
	require.NoError(t, err, "error making HTTP request")

	assert.Contains(t, resp.Header().Get("Set-Cookie"), jwtToken, "response cookie is unset")
}

func TestJwtAuthorizationMiddlewareNoCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockUsersService(ctrl)

	r := chi.NewRouter()
	r.Get("/", NewJwtHandler(mockHandler, service))

	ts := httptest.NewServer(r)
	defer ts.Close()

	userID := "test_user_id"

	req := resty.New().R()
	req.Method = "GET"
	req.URL = ts.URL

	service.EXPECT().Create(gomock.Any()).Return(&model.User{
		UserID: userID,
	}, nil)

	resp, err := req.Send()
	require.NoError(t, err, "error making HTTP request")

	assert.Contains(t, resp.Header().Get("Set-Cookie"), "Authorization=", "response cookie is unset")
}

func TestJwtAuthorizationMiddlewareNoUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockUsersService(ctrl)

	r := chi.NewRouter()
	r.Get("/", NewJwtHandler(mockHandler, service))

	ts := httptest.NewServer(r)
	defer ts.Close()

	userID := "test_user_id"
	jwtToken, err := buildJWT(userID)
	assert.NoError(t, err, "error while building jwt")

	req := resty.New().R()
	req.Method = "GET"
	req.URL = ts.URL
	req.Cookies = []*http.Cookie{
		{
			Name:  authCookieName,
			Value: jwtToken,
		},
	}

	service.EXPECT().GetByID(gomock.Any(), userID).Return(nil, model.ErrDBNotFound)

	resp, err := req.Send()
	require.NoError(t, err, "error making HTTP request")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode(), "response cookie is unset")
}
