package service

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type Storage interface {
	SaveURL(ctx context.Context, url string, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
}

type UserData struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

type URLService struct {
	storage   Storage
	userMutex sync.Mutex
	userUrls  map[string][]UserData
}

func NewURLService(storage Storage) *URLService {
	return &URLService{
		storage:   storage,
		userUrls:  make(map[string][]UserData),
		userMutex: sync.Mutex{},
	}
}

func (service *URLService) SaveURL(ctx context.Context, url string, userID string, baseURL string) (string, error) {
	key := uuid.New().String()
	err := service.storage.SaveURL(ctx, url, key)

	service.userMutex.Lock()
	service.userUrls[userID] = append(service.userUrls[userID], UserData{
		FullURL:  url,
		ShortURL: baseURL + "/" + key,
	})
	service.userMutex.Unlock()

	if err != nil {
		return "", err
	}

	return key, nil
}

func (service *URLService) GetURL(ctx context.Context, url string) (string, error) {
	return service.storage.GetURL(ctx, url)
}

func (service *URLService) GetUserData(userID string) []UserData {
	return service.userUrls[userID]
}
