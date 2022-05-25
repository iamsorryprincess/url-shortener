package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
)

type UserData struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

type URLService struct {
	storage   storage.Storage
	userMutex sync.Mutex
	userUrls  map[string][]UserData
}

func NewURLService(storage storage.Storage) *URLService {
	return &URLService{
		storage:   storage,
		userUrls:  make(map[string][]UserData),
		userMutex: sync.Mutex{},
	}
}

type URLUniqueError struct {
	OriginalURL string
	ShortURL    string
	err         error
}

func (e *URLUniqueError) Error() string {
	return fmt.Sprintf("original url %s already exist; shortURL: %s", e.OriginalURL, e.ShortURL)
}

func (e *URLUniqueError) Unwrap() error {
	return e.err
}

func (service *URLService) SaveURL(ctx context.Context, url string, userID string, baseURL string) (string, error) {
	key := uuid.New().String()
	err := service.storage.SaveURL(ctx, url, key)

	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExist) {
			shortURL, getErr := service.storage.GetByOriginalURL(ctx, url)

			if getErr != nil {
				return "", getErr
			}

			return "", &URLUniqueError{
				OriginalURL: url,
				ShortURL:    shortURL,
			}
		}

		return "", err
	}

	service.userMutex.Lock()
	service.userUrls[userID] = append(service.userUrls[userID], UserData{
		FullURL:  url,
		ShortURL: baseURL + "/" + key,
	})
	service.userMutex.Unlock()
	return key, nil
}

func (service *URLService) GetURL(ctx context.Context, url string) (string, error) {
	return service.storage.GetURL(ctx, url)
}

func (service *URLService) GetUserData(userID string) []UserData {
	return service.userUrls[userID]
}

type URLInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type URLResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (service *URLService) SaveBatch(ctx context.Context, baseURL string, input []URLInput) ([]URLResult, error) {
	batchData := make([]storage.URLInput, len(input))
	result := make([]URLResult, len(input))

	for index, inputData := range input {
		id := uuid.New().String()
		batchData[index] = storage.URLInput{
			ShortURL: id,
			FullURL:  inputData.OriginalURL,
		}
		result[index] = URLResult{
			CorrelationID: inputData.CorrelationID,
			ShortURL:      baseURL + "/" + id,
		}
	}

	err := service.storage.SaveBatch(ctx, batchData)

	if err != nil {
		return nil, err
	}

	return result, nil
}
