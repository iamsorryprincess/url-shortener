package app

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/handlers"
	"github.com/iamsorryprincess/url-shortener/internal/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/service"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
	"github.com/iamsorryprincess/url-shortener/pkg/hash"
	_ "github.com/jackc/pgx/stdlib"
)

type Configuration struct {
	Address            string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL            string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StoragePath        string `env:"FILE_STORAGE_PATH"`
	DbConnectionString string `env:"DATABASE_DSN" envDefault:""`
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
	flag.StringVar(&configuration.DbConnectionString, "d", configuration.DbConnectionString, "db connection string")
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

	keyManager, err := hash.NewGcmKeyManager()

	if err != nil {
		log.Fatal(err)
		return
	}

	if configuration.DbConnectionString == "" {
		log.Fatal("empty db connection string")
		return
	}

	db, err := sql.Open("pgx", configuration.DbConnectionString)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
		return
	}

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Gzip)
	r.Use(middleware.Cookie(keyManager))

	r.Post("/", handlers.RawMakeShortURLHandler(urlService, configuration.BaseURL))
	r.Post("/api/shorten", handlers.JSONMakeShortURLHandler(urlService, configuration.BaseURL))
	r.Get("/{URL}", handlers.GetFullURLHandler(urlService))
	r.Get("/api/user/urls", handlers.GetUserUrls(urlService))

	r.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		err := db.PingContext(request.Context())

		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(configuration.Address, r))
}
