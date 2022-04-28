package infrastructure

import (
	"cmd/shortener/main.go/internal/service"
	"sync"
)

type Storage struct {
	mutex    sync.Mutex
	localMap map[string]string
}

func (storage *Storage) SaveURL(url string, shortURL string) {
	storage.mutex.Lock()
	storage.localMap[shortURL] = url
	storage.mutex.Unlock()
}

func (storage *Storage) GetURL(shortURL string) string {
	return storage.localMap[shortURL]
}

func InitStorage() service.Storage {
	return &Storage{
		localMap: make(map[string]string),
		mutex:    sync.Mutex{},
	}
}
