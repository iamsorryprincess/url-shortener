package handlers

import (
	"cmd/shortener/main.go/internal/service"
	"io"
	"net/http"
)

func PostURL(writer http.ResponseWriter, request *http.Request, urlService *service.URLService) {
	bytes, readErr := io.ReadAll(request.Body)

	if readErr != nil {
		http.Error(writer, readErr.Error(), http.StatusInternalServerError)
		return
	}

	if len(bytes) == 0 {
		http.Error(writer, "empty body", http.StatusBadRequest)
		return
	}

	url := string(bytes)
	shortURL := urlService.SaveURL(url)
	writer.WriteHeader(http.StatusCreated)
	_, err := writer.Write([]byte("http://localhost:8080/" + shortURL))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetURL(writer http.ResponseWriter, request *http.Request, urlService *service.URLService) {
	url := request.RequestURI[1:len(request.RequestURI)]

	if url == "" {
		http.Error(writer, "empty url", http.StatusBadRequest)
		return
	}

	targetURL := urlService.GetURL(url)

	if targetURL == "" {
		http.NotFound(writer, request)
		return
	}

	writer.Header().Set("Location", targetURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}
