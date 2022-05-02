package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/handlers"
	"github.com/iamsorryprincess/url-shortener/internal/service"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
)

type Configuration struct {
	Address string
	BaseURL string
}

func parseConfiguration() *Configuration {
	return &Configuration{
		Address: os.Getenv("SERVER_ADDRESS"),
		BaseURL: os.Getenv("BASE_URL"),
	}
}

func Run() {
	configuration := parseConfiguration()
	urlStorage := storage.InitInMemoryStorage()
	urlService := service.InitURLService(urlStorage)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", handlers.RawMakeShortURLHandler(urlService, configuration.BaseURL))
	r.Post("/api/shorten", handlers.JSONMakeShortURLHandler(urlService, configuration.BaseURL))
	r.Get("/{URL}", handlers.GetFullURLHandler(urlService))

	err := http.ListenAndServe(fmt.Sprintf("localhost%s", configuration.Address), r)

	if err != nil {
		log.Fatal(err)
	}
}
