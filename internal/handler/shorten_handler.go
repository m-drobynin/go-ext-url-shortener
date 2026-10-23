package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func handleCreatedJSONResponse(config SaveHandlerConfig, w http.ResponseWriter, code string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var response = ShortenResponse{
		Result: config.GetBaseURL() + "/" + code,
	}

	res, err := json.Marshal(response)

	if err != nil {
		return err
	}

	w.Write(res)
	return nil
}

func BuildShortenHandler(config SaveHandlerConfig, service service.URLService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			handleError(w, ErrBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			handleError(w, fmt.Errorf("error while reading body: %v: %w", err, ErrInternal))
			return
		}

		if len(body) == 0 {
			handleError(w, fmt.Errorf("body is empty: %v: %w", err, ErrBadRequest))
			return
		}

		var request ShortenRequest

		if err = json.Unmarshal(body, &request); err != nil {
			handleError(w, fmt.Errorf("error while unmarshalling: %v: %w", err, ErrInternal))
			return
		}

		if len(request.URL) == 0 {
			handleError(w, fmt.Errorf("url is empty: %v: %w", err, ErrBadRequest))
			return
		}

		code, err := service.SaveURL(request.URL)

		if err != nil {
			handleError(w, fmt.Errorf("error while saving url: %v: %w", err, ErrInternal))
			return
		}

		err = handleCreatedJSONResponse(config, w, *code)

		if err != nil {
			handleError(w, fmt.Errorf("error while response marshalling: %v: %w", err, ErrInternal))
			return
		}
	}
}

// Добавьте в код сервера новый эндпоинт POST /api/shorten, который будет принимать в теле запроса JSON-объект {"url":"<some_url>"} и возвращать в ответ объект {"result":"<short_url>"}.
// Запрос может иметь такой вид:

// POST http://localhost:8080/api/shorten HTTP/1.1
// Host: localhost:8080
// Content-Type: application/json
// {
//   "url": "https://practicum.yandex.ru"
// }
// Ответ может быть таким:

// HTTP/1.1 201 OK
// Content-Type: application/json
// Content-Length: 30
// {
//  "result": "http://localhost:8080/EwHXdJfB"
// }
// Удостоверьтесь, что в ответе от сервера присутствует HTTP-заголовок Content-Type со значением application/json. Он указывает клиенту, в каком формате передано тело ответа.
// Также не забудьте добавить тесты на новый эндпоинт, как и на предыдущие.
