package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
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
				handler:     JSONMakeShortURLHandler(urlService),
				t:           t,
			}
			testHandler(testInfo)
		})
	}
}

func TestRawMakeShortURLHandler(t *testing.T) {
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
			testInfo := TestHandler{
				url:         "/",
				method:      http.MethodPost,
				query:       "",
				body:        test.body,
				statusCode:  test.expectedStatusCode,
				header:      "Content-Type",
				headerValue: test.expectedContentType,
				handler:     RawMakeShortURLHandler(urlService),
				t:           t,
			}
			testHandler(testInfo)
		})
	}
}

func TestGetFullURLHandler(t *testing.T) {
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
