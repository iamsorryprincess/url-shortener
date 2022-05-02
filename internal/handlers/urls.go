package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iamsorryprincess/url-shortener/internal/service"
)

func RawMakeShortURLHandler(urlService *service.URLService, baseURL string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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
		shortenURL := urlService.SaveURL(url)
		writer.WriteHeader(http.StatusCreated)
		_, err := writer.Write([]byte(fmt.Sprintf("%s/%s", baseURL, shortenURL)))

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	Result string `json:"result"`
}

func JSONMakeShortURLHandler(urlService *service.URLService, baseURL string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bytes, readErr := io.ReadAll(request.Body)

		if readErr != nil {
			http.Error(writer, readErr.Error(), http.StatusInternalServerError)
			return
		}

		if request.Header.Get("Content-Type") != "application/json" {
			http.Error(writer, "invalid content type", 400)
			return
		}

		if len(bytes) == 0 {
			http.Error(writer, "empty body", http.StatusBadRequest)
			return
		}

		reqBody := URLRequest{}
		deserializeErr := json.Unmarshal(bytes, &reqBody)

		if deserializeErr != nil {
			http.Error(writer, deserializeErr.Error(), http.StatusInternalServerError)
			return
		}

		if reqBody.URL == "" {
			http.Error(writer, "url is empty", http.StatusBadRequest)
			return
		}

		shortenURL := urlService.SaveURL(reqBody.URL)
		response := URLResponse{
			Result: fmt.Sprintf("%s/%s", baseURL, shortenURL),
		}

		responseBytes, serializeErr := json.Marshal(&response)

		if serializeErr != nil {
			http.Error(writer, serializeErr.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(responseBytes)
	}
}

func GetFullURLHandler(urlService *service.URLService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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
}
