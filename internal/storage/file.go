package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"sync"
)

type FileStorage struct {
	mutex   sync.Mutex
	storage map[string]string
	file    *os.File
	encoder *json.Encoder
}

type storageData struct {
	ShortURL string `json:"shortUrl"`
	FullURL  string `json:"fullUrl"`
}

func NewFileStorage(filepath string) (*FileStorage, error) {
	file, openFileErr := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if openFileErr != nil {
		return nil, openFileErr
	}

	storage := make(map[string]string)
	reader := bufio.NewReader(file)

	for {
		bytes, readErr := reader.ReadBytes('\n')

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return nil, readErr
		}

		data := &storageData{}
		unmarshalErr := json.Unmarshal(bytes, data)

		if unmarshalErr != nil {
			return nil, unmarshalErr
		}

		storage[data.ShortURL] = data.FullURL
	}

	return &FileStorage{
		mutex:   sync.Mutex{},
		storage: storage,
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (s *FileStorage) SaveURL(url string, shortURL string) error {
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

func (s *FileStorage) GetURL(shortURL string) string {
	return s.storage[shortURL]
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}
