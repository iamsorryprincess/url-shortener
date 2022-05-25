package storage

import (
	"context"
	"errors"
	"sync"
)

type inMemoryStorage struct {
	mutex    sync.Mutex
	localMap map[string]string
}

func NewInMemoryStorage() Storage {
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

func (storage *inMemoryStorage) SaveBatch(ctx context.Context, batchData []URLInput) error {
	return errors.New("method not implemented")
}

func (storage *inMemoryStorage) GetByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	return "", errors.New("method not implemented")
}
