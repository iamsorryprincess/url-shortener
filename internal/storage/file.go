package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/iamsorryprincess/url-shortener/internal/service"
)

type fileStorage struct {
	mutex   sync.Mutex
	storage map[string]string
	file    *os.File
	encoder *json.Encoder
}

type storageData struct {
	ShortURL string `json:"shortUrl"`
	FullURL  string `json:"fullUrl"`
}

func NewFileStorage(filepath string) (service.Storage, io.Closer, error) {
	file, openFileErr := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if openFileErr != nil {
		return nil, nil, openFileErr
	}

	storage := make(map[string]string)
	reader := bufio.NewReader(file)

	for {
		bytes, readErr := reader.ReadBytes('\n')

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			file.Close()
			return nil, nil, readErr
		}

		data := &storageData{}
		unmarshalErr := json.Unmarshal(bytes, data)

		if unmarshalErr != nil {
			file.Close()
			return nil, nil, unmarshalErr
		}

		storage[data.ShortURL] = data.FullURL
	}

	return &fileStorage{
		mutex:   sync.Mutex{},
		storage: storage,
		file:    file,
		encoder: json.NewEncoder(file),
	}, file, nil
}

func (s *fileStorage) SaveURL(ctx context.Context, url string, shortURL string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.storage[shortURL] = url

	data := &storageData{
		ShortURL: shortURL,
		FullURL:  url,
	}

	err := s.encoder.Encode(data)

	if err != nil {
		return err
	}

	return nil
}

func (s *fileStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	return s.storage[shortURL], nil
}
