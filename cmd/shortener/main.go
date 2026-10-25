package main

import (
	"fmt"
	"log"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"m-drobynin/go-ext-url-shortener/internal/logger"
	"m-drobynin/go-ext-url-shortener/internal/middleware"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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

	log, err := logger.NewLogger("info")

	if err != nil {
		return err
	}

	// database := model.NewDatabase()
	fileStorage, err := model.NewFileStorage(config.FileStoragePath)

	if err != nil {
		return err
	}

	service := service.NewURLServiceImpl(fileStorage)
	handler := prepareApplication(config, service, log)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.NetAddress.Port), handler)
}

func prepareApplication(config *config.AppConfig, service service.URLService, log *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NewLoggerMiddleware(log), middleware.WithGzip)

	var saveHandler = handler.BuildSaveHandler(config, service)
	var getHandler = handler.BuildGetHandler(service)
	var shortenHandler = handler.BuildShortenHandler(config, service)

	r.Get("/{urlCode}", getHandler)
	r.Post("/", saveHandler)
	r.Post("/api/shorten", shortenHandler)

	return r
}
