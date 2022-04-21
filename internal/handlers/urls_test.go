package handlers

import (
	"bytes"
	"cmd/shortener/main.go/internal/service"
	"cmd/shortener/main.go/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostURL(t *testing.T) {
	urlStorage := storage.InitInMemoryStorage()
	urlService := service.InitURLService(urlStorage)

	tests := []struct {
		name                string
		body                string
		expectedStatusCode  int
		expectedContentType string
	}{
		{
			name:                "check with empty url",
			body:                "",
			expectedStatusCode:  400,
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name:                "check non empty url",
			body:                "https://www.youtube.com/",
			expectedStatusCode:  201,
			expectedContentType: "",
		},
		{
			name:                "check same url",
			body:                "https://www.youtube.com/",
			expectedStatusCode:  201,
			expectedContentType: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody := []byte(test.body)
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			writer := httptest.NewRecorder()

			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				PostURL(writer, request, urlService)
			})

			handler.ServeHTTP(writer, request)
			response := writer.Result()
			defer response.Body.Close()

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code: %d\nActual status code: %d\n", test.expectedStatusCode, response.StatusCode)
			}

			contentType := response.Header.Get("Content-Type")

			if contentType != test.expectedContentType {
				t.Errorf("Expected content type: %s\nActual content type: %s\n", test.expectedContentType, contentType)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	urlStorage := storage.InitInMemoryStorage()
	urlService := service.InitURLService(urlStorage)
	url := "https://www.youtube.com/"
	shortURL := urlService.SaveURL(url)

	tests := []struct {
		name               string
		query              string
		expectedStatusCode int
		locationHeader     string
	}{
		{
			name:               "test with empty url",
			query:              "",
			expectedStatusCode: 400,
			locationHeader:     "",
		},
		{
			name:               "test with not empty url",
			query:              shortURL,
			expectedStatusCode: 307,
			locationHeader:     url,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+test.query, nil)
			writer := httptest.NewRecorder()

			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				GetURL(writer, request, urlService)
			})

			handler.ServeHTTP(writer, request)
			response := writer.Result()
			defer response.Body.Close()

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code: %d\nActual status code: %d\n", test.expectedStatusCode, response.StatusCode)
			}

			locationHeader := response.Header.Get("Location")

			if locationHeader != test.locationHeader {
				t.Errorf("Expected location header: %s\nActual location header: %s\n", test.locationHeader, locationHeader)
			}
		})
	}
}
