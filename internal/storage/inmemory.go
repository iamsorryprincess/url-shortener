package storage

import (
	"context"
	"errors"
	"sync"
)

type inMemoryStorage struct {
	mutex    sync.Mutex
	urls     map[string]string
	userData map[string][]UserData
}

func NewInMemoryStorage() Storage {
	return &inMemoryStorage{
		urls:     make(map[string]string),
		userData: make(map[string][]UserData),
		mutex:    sync.Mutex{},
	}
}

func (storage *inMemoryStorage) SaveURL(ctx context.Context, input URLInput) error {
	storage.mutex.Lock()
	storage.urls[input.ShortURL] = input.FullURL
	storage.userData[input.UserID] = append(storage.userData[input.UserID], UserData{
		ShortURL: input.ShortURL,
		FullURL:  input.FullURL,
	})
	storage.mutex.Unlock()
	return nil
}

func (storage *inMemoryStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	return storage.urls[shortURL], nil
}

func (storage *inMemoryStorage) GetURLsByUserID(ctx context.Context, userID string) ([]UserData, error) {
	return storage.userData[userID], nil
}

func (storage *inMemoryStorage) SaveBatch(ctx context.Context, batchData []URLInput) error {
	return errors.New("method not implemented")
}

func (storage *inMemoryStorage) GetByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	return "", errors.New("method not implemented")
}
