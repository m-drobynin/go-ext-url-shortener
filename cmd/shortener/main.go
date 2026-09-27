package main

import (
	"fmt"
	"m-drobynin/go-ext-url-shortener/internal/config"
	"m-drobynin/go-ext-url-shortener/internal/handler"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	config := config.GetAppConfig()
	database := model.CreateDatabase()
	handler := prepareApplication(&config, &database)

	fmt.Println(&config)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.NetAddress.Port), handler)
}

func prepareApplication(config *config.AppConfig, database *model.Database) chi.Router {
	r := chi.NewRouter()

	var saveHandler = handler.BuildSaveHandler(config, database)
	var getHandler = handler.BuildGetHandler(config, database)

	r.Get("/{urlCode}", getHandler)
	r.Post("/", saveHandler)

	return r
}
