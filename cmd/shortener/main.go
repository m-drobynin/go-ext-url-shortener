package main

import (
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"
)

const urlLength = 10

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var database = model.CreateDatabase()
	var shortenerHandler = handler.BuildShortenerHandler(&database)

	http.HandleFunc("/", shortenerHandler)

	return http.ListenAndServe(`:8080`, nil)
}
