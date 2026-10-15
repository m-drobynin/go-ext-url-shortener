package handler

import (
	"io"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenHandler(t *testing.T) {
	config := config.GetAppConfigDefaults()
	service := &TestSaveURLServiceImpl{}

	handler := BuildShortenHandler(&config, service)

	jsonContentType := "application/json"
	emptyStr := ""
	validURL := "http://example.com/path"
	badRequestBody := http.StatusText(http.StatusBadRequest)
	validBodyMessage := `{"result":"` + config.GetBaseURL() + `/urlcodetest"}`

	cases := []SaveTestCase{
		{
			name: "empty url nil",

			url:         nil,
			contentType: &jsonContentType,

			expectedCode: http.StatusBadRequest,
			expectedBody: &badRequestBody,
		},
		{
			name: "empty url",

			url:         &emptyStr,
			contentType: &jsonContentType,

			expectedCode: http.StatusBadRequest,
			expectedBody: &badRequestBody,
		},
		{
			name: "invalid header",

			url: &validURL,

			expectedCode: http.StatusBadRequest,
			expectedBody: &badRequestBody,
		},
		{
			name:        "valid request",
			url:         &validURL,
			contentType: &jsonContentType,

			expectedCode: http.StatusCreated,
			expectedBody: &validBodyMessage,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var reader io.Reader = nil

			if tc.url != nil {
				req := `{"url":"` + *tc.url + `"}`
				reader = strings.NewReader(req)
			}

			r := httptest.NewRequest(http.MethodPost, "/", reader)

			if tc.contentType != nil {
				r.Header.Add("Content-Type", *tc.contentType)
			}

			w := httptest.NewRecorder()

			handler(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Regexp(t, *tc.expectedBody, w.Body.String(), "Невалидный ответ")
		})
	}
}
