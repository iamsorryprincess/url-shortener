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
		shortenURL, serviceErr := urlService.SaveURL(url, getUserId(request), baseURL)

		if serviceErr != nil {
			http.Error(writer, "internal error", http.StatusInternalServerError)
			return
		}

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

		shortenURL, serviceErr := urlService.SaveURL(reqBody.URL, getUserId(request), baseURL)

		if serviceErr != nil {
			http.Error(writer, "internal error", http.StatusInternalServerError)
			return
		}

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

func GetUserUrls(urlService *service.URLService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data := urlService.GetUserData(getUserId(request))

		if data == nil {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		bytes, err := json.Marshal(data)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)
	}
}

func getUserId(request *http.Request) string {
	value, ok := request.Context().Value("user_id").(string)

	if !ok {
		return ""
	}

	return value
}
