package main

import (
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"net/http"
)

const urlLength = 10

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	http.HandleFunc("/", handler.HandleShortenerRequest)

	return http.ListenAndServe(`:8080`, nil)
}
