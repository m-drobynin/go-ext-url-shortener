package main

import (
	"database/sql"
	"fmt"
	"log"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"m-drobynin/go-ext-url-shortener/internal/logger"
	"m-drobynin/go-ext-url-shortener/internal/middleware"
	"m-drobynin/go-ext-url-shortener/internal/repository"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"m-drobynin/go-ext-url-shortener/internal/storage"
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

	var db *sql.DB

	if config.DatabaseDSN != "" {
		db, err = repository.PrepareDBConnection(config.DatabaseDSN)

		if err != nil {
			return err
		}
	}

	urlsStorage, err := prepareUrlsStorage(config, db)

	if err != nil {
		return err
	}

	usersStorage, err := prepareUsersStorage(config, db)

	if err != nil {
		return err
	}

	urlsService := service.NewURLServiceImpl(urlsStorage)
	usersService := service.NewUsersServiceImpl(usersStorage)

	router := prepareApplication(config, urlsService, usersService, log)

	if db != nil {
		router.Get("/ping", handler.BuildPingHandler(db))
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", config.NetAddress.Port), router)
}

func prepareUrlsStorage(config *config.AppConfig, db *sql.DB) (storage.UrlsDBProvider, error) {
	if db != nil {
		return storage.NewUrlsPostgresDB(db)
	}

	if config.FileStoragePath != "" {
		return storage.NewUrlsFileStorage(config.FileStoragePath)
	}

	return storage.NewUrlsMemoryDatabase(), nil
}

func prepareUsersStorage(config *config.AppConfig, db *sql.DB) (storage.UsersDBProvider, error) {
	if db != nil {
		return storage.NewUsersPostgresDB(db)
	}

	if config.FileStoragePath != "" {
		return storage.NewUsersFileStorage(config.UsersFileStoragePath)
	}

	return storage.NewUsersMemoryDatabase(), nil
}

func prepareApplication(config *config.AppConfig, urlsService service.URLService, usersService service.UsersService, log *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NewLoggerMiddleware(log), middleware.WithGzip, middleware.WithJwtAuthorization)

	var saveHandler = handler.BuildSaveHandler(config, urlsService)
	var getHandler = handler.BuildGetHandler(urlsService)
	var shortenHandler = handler.BuildShortenHandler(config, urlsService)
	var batchShortenHandler = handler.BuildBatchShortenHandler(config, urlsService)

	r.Get("/{urlCode}", getHandler)
	r.Post("/", saveHandler)
	r.Post("/api/shorten", shortenHandler)
	r.Post("/api/shorten/batch", batchShortenHandler)

	return r
}
