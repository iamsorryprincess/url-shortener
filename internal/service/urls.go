package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
)

type URLService struct {
	storage   storage.Storage
	userMutex sync.Mutex
	baseURL   string
}

func NewURLService(storage storage.Storage, baseURL string) *URLService {
	return &URLService{
		storage:   storage,
		userMutex: sync.Mutex{},
		baseURL:   baseURL,
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

func (service *URLService) SaveURL(ctx context.Context, url string, userID string) (string, error) {
	key := uuid.New().String()
	err := service.storage.SaveURL(ctx, storage.URLInput{
		FullURL:  url,
		ShortURL: key,
		UserID:   userID,
	})

	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExist) {
			shortURL, getErr := service.storage.GetByOriginalURL(ctx, url)

			if getErr != nil {
				return "", getErr
			}

			return "", &URLUniqueError{
				OriginalURL: url,
				ShortURL:    service.baseURL + "/" + shortURL,
			}
		}

		return "", err
	}

	return service.baseURL + "/" + key, nil
}

func (service *URLService) GetURL(ctx context.Context, url string) (string, error) {
	return service.storage.GetURL(ctx, url)
}

func (service *URLService) GetUserData(ctx context.Context, userID string) ([]storage.UserData, error) {
	result, err := service.storage.GetURLsByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	for index, item := range result {
		result[index].ShortURL = service.baseURL + "/" + item.ShortURL
	}

	return result, nil
}

type URLInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type URLResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (service *URLService) SaveBatch(ctx context.Context, input []URLInput) ([]URLResult, error) {
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
			ShortURL:      service.baseURL + "/" + id,
		}
	}

	err := service.storage.SaveBatch(ctx, batchData)

	if err != nil {
		return nil, err
	}

	return result, nil
}
