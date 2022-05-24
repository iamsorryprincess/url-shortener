package storage

import (
	"context"
	"sync"

	"github.com/iamsorryprincess/url-shortener/internal/service"
)

type inMemoryStorage struct {
	mutex    sync.Mutex
	localMap map[string]string
}

func NewInMemoryStorage() service.Storage {
	return &inMemoryStorage{
		localMap: make(map[string]string),
		mutex:    sync.Mutex{},
	}
}

func (storage *inMemoryStorage) SaveURL(ctx context.Context, url string, shortURL string) error {
	storage.mutex.Lock()
	storage.localMap[shortURL] = url
	storage.mutex.Unlock()
	return nil
}

func (storage *inMemoryStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	return storage.localMap[shortURL], nil
}
