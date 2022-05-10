package app

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/handlers"
	"github.com/iamsorryprincess/url-shortener/internal/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/service"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
)

type Configuration struct {
	Address     string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StoragePath string `env:"FILE_STORAGE_PATH"`
}

func parseConfiguration() (*Configuration, error) {
	configuration := &Configuration{}
	err := env.Parse(configuration)

	if err != nil {
		return nil, err
	}

	flag.StringVar(&configuration.Address, "a", configuration.Address, "server address")
	flag.StringVar(&configuration.BaseURL, "b", configuration.BaseURL, "base url")
	flag.StringVar(&configuration.StoragePath, "f", configuration.StoragePath, "file storage path")
	flag.Parse()
	return configuration, nil
}

func Run() {
	configuration, confErr := parseConfiguration()

	if confErr != nil {
		log.Fatal(confErr)
		return
	}

	var urlService *service.URLService

	if configuration.StoragePath == "" {
		inMemoryStorage := storage.NewInMemoryStorage()
		urlService = service.NewURLService(inMemoryStorage)
	} else {
		fileStorage, err := storage.NewFileStorage(configuration.StoragePath)

		if err != nil {
			log.Fatal(err)
			return
		}

		urlService = service.NewURLService(fileStorage)
	}

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Gzip)

	r.Post("/", handlers.RawMakeShortURLHandler(urlService, configuration.BaseURL))
	r.Post("/api/shorten", handlers.JSONMakeShortURLHandler(urlService, configuration.BaseURL))
	r.Get("/{URL}", handlers.GetFullURLHandler(urlService))

	err := http.ListenAndServe(configuration.Address, r)

	if err != nil {
		log.Fatal(err)
	}
}
