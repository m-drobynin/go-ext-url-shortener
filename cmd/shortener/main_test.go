package main

import (
	"io"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

func testRequest(t *testing.T, method,
	path string, reqBody *string) (*resty.Response, string) {
	client := resty.New()
	client.SetRedirectPolicy(resty.RedirectNoPolicy())

	req := client.R()
	req.Method = method
	req.URL = path

	if reqBody != nil {
		req.SetContentType("text/plain")
		req.Body = *reqBody
	}

	resp, err := req.Send()
	require.NoError(t, err, "error making HTTP request")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "error reading response body")

	return resp, string(body)
}

func TestRouterHappyRoute(t *testing.T) {
	var database = model.CreateDatabase()
	handler := buildRouter(&database)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	originalURL := "https://example.com"
	localhost := "http://localhost:8080/"

	saveResponse, saveBody := testRequest(t, "POST", ts.URL+"/", &originalURL)

	assert.Equal(t, 201, saveResponse.StatusCode())
	assert.Regexp(t, localhost, saveBody)

	urlCode := strings.Replace(saveBody, localhost, "", 1)

	getResponse, _ := testRequest(t, "GET", ts.URL+"/"+urlCode, nil)

	assert.Equal(t, getResponse.Header().Get("Location"), originalURL)
	assert.Equal(t, http.StatusTemporaryRedirect, getResponse.StatusCode())
}
