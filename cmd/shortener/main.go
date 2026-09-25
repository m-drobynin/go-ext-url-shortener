package main

import (
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const urlLength = 10

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var database = model.CreateDatabase()
	handler := buildRouter(&database)

	return http.ListenAndServe(`:8080`, handler)
}

func buildRouter(database *model.Database) chi.Router {
	r := chi.NewRouter()

	var saveHandler = handler.BuildSaveHandler(database)
	var getHandler = handler.BuildGetHandler(database)

	r.Get("/{urlCode}", getHandler)
	r.Post("/", saveHandler)

	return r
}
