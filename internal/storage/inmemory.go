package storage

import (
	"sync"
)

type InMemoryStorage struct {
	mutex    sync.Mutex
	localMap map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		localMap: make(map[string]string),
		mutex:    sync.Mutex{},
	}
}

func (storage *InMemoryStorage) SaveURL(url string, shortURL string) error {
	storage.mutex.Lock()
	storage.localMap[shortURL] = url
	storage.mutex.Unlock()
	return nil
}

func (storage *InMemoryStorage) GetURL(shortURL string) string {
	return storage.localMap[shortURL]
}
