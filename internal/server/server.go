package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/config"
	"github.com/iamsorryprincess/url-shortener/internal/handlers"
	"github.com/iamsorryprincess/url-shortener/internal/middleware"
	"github.com/iamsorryprincess/url-shortener/internal/service"
	"github.com/iamsorryprincess/url-shortener/pkg/hash"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(
	configuration *config.Configuration,
	service *service.URLService,
	keyManager hash.KeyManager,
	db *sql.DB) *Server {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Gzip)
	r.Use(middleware.Cookie(keyManager))

	r.Post("/", handlers.RawMakeShortURLHandler(service))
	r.Post("/api/shorten", handlers.JSONMakeShortURLHandler(service))
	r.Post("/api/shorten/batch", handlers.SaveBatchURLHandler(service))
	r.Get("/{URL}", handlers.GetFullURLHandler(service))
	r.Get("/api/user/urls", handlers.GetUserUrls(service))
	r.Delete("/api/user/urls", handlers.DeleteBatchURLHandler(service))

	if configuration.DBConnectionString != "" {
		r.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
			pingErr := db.PingContext(request.Context())

			if pingErr != nil {
				log.Println(pingErr)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			writer.WriteHeader(http.StatusOK)
		})
	}

	return &Server{
		httpServer: &http.Server{
			Addr:    configuration.Address,
			Handler: r,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
