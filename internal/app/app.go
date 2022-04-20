package app

import (
	"cmd/shortener/main.go/internal/infrastructure"
	"cmd/shortener/main.go/internal/service"
	"io"
	"log"
	"net/http"
)

func postURL(writer http.ResponseWriter, request *http.Request, urlService *service.URLService) {
	bytes, readErr := io.ReadAll(request.Body)

	if readErr != nil {
		http.Error(writer, readErr.Error(), http.StatusInternalServerError)
		return
	}

	if len(bytes) == 0 {
		http.Error(writer, "empty body", http.StatusBadRequest)
		return
	}

	var url = string(bytes)
	var shortURL = urlService.SaveURL(url)
	writer.WriteHeader(http.StatusCreated)
	_, err := writer.Write([]byte("http://localhost:8080/" + shortURL))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getURL(writer http.ResponseWriter, request *http.Request, urlService *service.URLService) {
	var url = request.RequestURI[1:len(request.RequestURI)]

	if url == "" {
		http.Error(writer, "empty url", http.StatusBadRequest)
		return
	}

	var targetURL = urlService.GetURL(url)

	if targetURL == "" {
		http.NotFound(writer, request)
		return
	}

	writer.Header().Set("Location", targetURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func Run() {
	var storage = infrastructure.InitStorage()
	var urlService = service.InitURLService(&storage)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			postURL(writer, request, urlService)
		case http.MethodGet:
			getURL(writer, request, urlService)
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
