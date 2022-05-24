package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iamsorryprincess/url-shortener/internal/service"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
)

type TestHandler struct {
	url         string
	method      string
	contentType string
	query       string
	body        string
	statusCode  int
	header      string
	headerValue string
	handler     http.HandlerFunc
	t           *testing.T
}

func testHandler(handlerTestInfo TestHandler) {
	reqBody := []byte(handlerTestInfo.body)
	request := httptest.NewRequest(handlerTestInfo.method, fmt.Sprintf("%s%s",
		handlerTestInfo.url,
		handlerTestInfo.query), bytes.NewBuffer(reqBody))

	if handlerTestInfo.contentType != "" {
		request.Header.Set("Content-Type", handlerTestInfo.contentType)
	}

	writer := httptest.NewRecorder()

	handlerTestInfo.handler.ServeHTTP(writer, request)
	response := writer.Result()
	defer response.Body.Close()

	if response.StatusCode != handlerTestInfo.statusCode {
		handlerTestInfo.t.Errorf("Expected status code: %d\nActual status code: %d\n", handlerTestInfo.statusCode, response.StatusCode)
	}

	contentType := response.Header.Get(handlerTestInfo.header)

	if contentType != handlerTestInfo.headerValue {
		handlerTestInfo.t.Errorf("Expected content type: %s\nActual content type: %s\n", handlerTestInfo.headerValue, contentType)
	}
}

func TestJSONMakeShortURLHandler(t *testing.T) {
	urlStorage := storage.NewInMemoryStorage()
	urlService := service.NewURLService(urlStorage)

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
			name:                "check with empty url in json body",
			body:                `{"url": ""}`,
			expectedStatusCode:  400,
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name:                "check non empty url",
			body:                `{"url": "https://www.youtube.com/"}`,
			expectedStatusCode:  201,
			expectedContentType: "application/json",
		},
		{
			name:                "check same url",
			body:                `{"url": "https://www.youtube.com/"}`,
			expectedStatusCode:  201,
			expectedContentType: "application/json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testInfo := TestHandler{
				url:         "/api/shorten",
				method:      http.MethodPost,
				contentType: "application/json",
				query:       "",
				body:        test.body,
				statusCode:  test.expectedStatusCode,
				header:      "Content-Type",
				headerValue: test.expectedContentType,
				handler:     JSONMakeShortURLHandler(urlService, "http://localhost:8080"),
				t:           t,
			}
			testHandler(testInfo)
		})
	}
}

func TestRawMakeShortURLHandler(t *testing.T) {
	urlStorage := storage.NewInMemoryStorage()
	urlService := service.NewURLService(urlStorage)

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
			testInfo := TestHandler{
				url:         "/",
				method:      http.MethodPost,
				query:       "",
				body:        test.body,
				statusCode:  test.expectedStatusCode,
				header:      "Content-Type",
				headerValue: test.expectedContentType,
				handler:     RawMakeShortURLHandler(urlService, "http://localhost:8080"),
				t:           t,
			}
			testHandler(testInfo)
		})
	}
}

func TestGetFullURLHandler(t *testing.T) {
	urlStorage := storage.NewInMemoryStorage()
	urlService := service.NewURLService(urlStorage)
	url := "https://www.youtube.com/"
	shortURL, err := urlService.SaveURL(context.Background(), url, "test", "http://localhost:8080")

	if err != nil {
		t.Fatal(err)
	}

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
			testInfo := TestHandler{
				url:         "/",
				method:      http.MethodGet,
				query:       test.query,
				body:        "",
				statusCode:  test.expectedStatusCode,
				header:      "Location",
				headerValue: test.locationHeader,
				handler:     GetFullURLHandler(urlService),
				t:           t,
			}
			testHandler(testInfo)
		})
	}
}

func TestGzipMiddleware(t *testing.T) {
	reqBody1 := URLRequest{
		URL: "https://practicum.yandex.ru/learn/go-developer/courses/d4f7d31d-bdf2-4d55-9845-3eb6d29448ea/sprints/21257/topics/2577c77d-8dac-4732-9d2d-48d0f9dbd57b/lessons/54bfce18-6b0e-4de4-a5bd-8a535196dff1/",
	}

	data, err := json.Marshal(reqBody1)

	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	urlStorage := storage.NewInMemoryStorage()
	urlService := service.NewURLService(urlStorage)
	handler := JSONMakeShortURLHandler(urlService, "http://localhost:8080")
	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, request)
	response := writer.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Errorf("invalid status code\nexpected: 201\nactual: %d", response.StatusCode)
	}

	if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
		t.Error("Content-Encoding header must not contain gzip")
	}
}
