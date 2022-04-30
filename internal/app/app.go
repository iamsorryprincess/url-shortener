package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/handlers"
	"github.com/iamsorryprincess/url-shortener/internal/service"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
)

func Run() {
	urlStorage := storage.InitInMemoryStorage()
	urlService := service.InitURLService(urlStorage)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", handlers.RawMakeShortURLHandler(urlService))
	r.Post("/api/shorten", handlers.JSONMakeShortURLHandler(urlService))
	r.Get("/{URL}", handlers.GetFullURLHandler(urlService))

	err := http.ListenAndServe("localhost:8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
