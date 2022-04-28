package app

import (
	"cmd/shortener/main.go/internal/handlers"
	"cmd/shortener/main.go/internal/service"
	"cmd/shortener/main.go/internal/storage"
	"log"
	"net/http"
)

func Run() {
	var urlStorage = storage.InitInMemoryStorage()
	var urlService = service.InitURLService(urlStorage)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			handlers.PostURL(writer, request, urlService)
		case http.MethodGet:
			handlers.GetURL(writer, request, urlService)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe("localhost:8080", nil)

	if err != nil {
		log.Fatal(err)
		return
	}
}
