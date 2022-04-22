package app

import (
	"cmd/shortener/main.go/internal/handlers"
	"cmd/shortener/main.go/internal/service"
	"cmd/shortener/main.go/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func Run() {
	urlStorage := storage.InitInMemoryStorage()
	urlService := service.InitURLService(urlStorage)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	postURLHandler := handlers.PostURLHandler(urlService)
	getURLHandler := handlers.GetURLHandler(urlService)

	r.Post("/", postURLHandler)
	r.Get("/{URL}", getURLHandler)

	err := http.ListenAndServe("localhost:8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
