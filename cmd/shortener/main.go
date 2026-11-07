package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"m-drobynin/go-ext-url-shortener/internal/logger"
	"m-drobynin/go-ext-url-shortener/internal/middleware"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func BuildPingHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := db.PingContext(ctx)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
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

	db, err := sql.Open("pgx", config.DatabaseDSN)

	if err != nil {
		return err
	}

	service := service.NewURLServiceImpl(fileStorage)
	router := prepareApplication(config, service, log)

	router.Get("/ping", BuildPingHandler(db))

	return http.ListenAndServe(fmt.Sprintf(":%d", config.NetAddress.Port), router)
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
