package main

import (
	"fmt"
	"log"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config := config.GetAppConfig()
	database := model.NewDatabase()
	service := service.NewURLServiceImpl(database)
	handler := prepareApplication(&config, service)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.NetAddress.Port), handler)
}

func prepareApplication(config *config.AppConfig, service service.URLService) chi.Router {
	r := chi.NewRouter()

	var saveHandler = handler.BuildSaveHandler(config, service)
	var getHandler = handler.BuildGetHandler(service)

	r.Get("/{urlCode}", getHandler)
	r.Post("/", saveHandler)

	return r
}
