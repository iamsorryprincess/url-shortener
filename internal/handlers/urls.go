package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/iamsorryprincess/url-shortener/internal/middleware"
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
		shortenURL, serviceErr := urlService.SaveURL(request.Context(), url, getUserID(request), baseURL)

		if serviceErr != nil {
			var urlErr *service.UrlUniqueError
			if errors.As(serviceErr, &urlErr) {
				writeURLResponseRaw(writer, baseURL, urlErr.ShortURL, http.StatusConflict)
				return
			}

			http.Error(writer, "internal error", http.StatusInternalServerError)
			return
		}

		writeURLResponseRaw(writer, baseURL, shortenURL, http.StatusCreated)
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

		shortenURL, serviceErr := urlService.SaveURL(request.Context(), reqBody.URL, getUserID(request), baseURL)

		if serviceErr != nil {
			var urlErr *service.UrlUniqueError
			if errors.As(serviceErr, &urlErr) {
				writeURLResponseJSON(writer, baseURL, urlErr.ShortURL, http.StatusConflict)
				return
			}

			http.Error(writer, "internal error", http.StatusInternalServerError)
			return
		}

		writeURLResponseJSON(writer, baseURL, shortenURL, http.StatusCreated)
	}
}

func GetFullURLHandler(urlService *service.URLService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		url := request.RequestURI[1:len(request.RequestURI)]

		if url == "" {
			http.Error(writer, "empty url", http.StatusBadRequest)
			return
		}

		targetURL, err := urlService.GetURL(request.Context(), url)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

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
		data := urlService.GetUserData(getUserID(request))

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

func SaveBatchURLHandler(urlService *service.URLService, baseURL string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Content-Type") != "application/json" {
			writer.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		bytes, err := io.ReadAll(request.Body)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(bytes) == 0 {
			http.Error(writer, "empty body", http.StatusBadRequest)
			return
		}

		var reqBody []service.URLInput
		err = json.Unmarshal(bytes, &reqBody)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		batchResult, err := urlService.SaveBatch(request.Context(), baseURL, reqBody)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(batchResult)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(result)
	}
}

func getUserID(request *http.Request) string {
	value, ok := request.Context().Value(middleware.CookieKey).(middleware.UserData)

	if !ok {
		return ""
	}

	return value.ID
}

func writeURLResponseRaw(writer http.ResponseWriter, baseURL string, shortenURL string, statusCode int) {
	writer.WriteHeader(statusCode)
	_, err := writer.Write([]byte(fmt.Sprintf("%s/%s", baseURL, shortenURL)))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeURLResponseJSON(writer http.ResponseWriter, baseURL string, shortenURL string, statusCode int) {
	response := URLResponse{
		Result: fmt.Sprintf("%s/%s", baseURL, shortenURL),
	}

	responseBytes, serializeErr := json.Marshal(&response)

	if serializeErr != nil {
		http.Error(writer, serializeErr.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(responseBytes)
}
