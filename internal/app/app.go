package app

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/handlers"
	"github.com/iamsorryprincess/url-shortener/internal/service"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
)

type Configuration struct {
	Address string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func parseConfiguration() (*Configuration, error) {
	configuration := &Configuration{}
	err := env.Parse(configuration)

	if err != nil {
		return nil, err
	}

	return configuration, nil
}

func Run() {
	configuration, confErr := parseConfiguration()

	if confErr != nil {
		log.Fatal(confErr)
		return
	}

	urlStorage := storage.InitInMemoryStorage()
	urlService := service.InitURLService(urlStorage)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", handlers.RawMakeShortURLHandler(urlService, configuration.BaseURL))
	r.Post("/api/shorten", handlers.JSONMakeShortURLHandler(urlService, configuration.BaseURL))
	r.Get("/{URL}", handlers.GetFullURLHandler(urlService))

	err := http.ListenAndServe(configuration.Address, r)

	if err != nil {
		log.Fatal(err)
	}
}
