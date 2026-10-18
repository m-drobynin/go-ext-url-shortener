package main

import (
	"fmt"
	"log"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"m-drobynin/go-ext-url-shortener/internal/logger"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"m-drobynin/go-ext-url-shortener/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := config.GetAppConfig()

	if err != nil {
		return err
	}

	if err = logger.Initialize("info"); err != nil {
		return err
	}

	database := model.NewDatabase()
	service := service.NewURLServiceImpl(database)
	handler := prepareApplication(config, service)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.NetAddress.Port), handler)
}

func applyMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return logger.WithLogging(utils.WithGzip(h))
}

func prepareApplication(config *config.AppConfig, service service.URLService) chi.Router {
	r := chi.NewRouter()

	var saveHandler = applyMiddleware(handler.BuildSaveHandler(config, service))
	var getHandler = applyMiddleware(handler.BuildGetHandler(service))
	var shortenHandler = applyMiddleware(handler.BuildShortenHandler(config, service))

	r.Get("/{urlCode}", getHandler)
	r.Post("/", saveHandler)
	r.Post("/api/shorten", shortenHandler)

	return r
}
