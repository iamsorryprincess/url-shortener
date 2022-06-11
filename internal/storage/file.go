package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

type fileStorage struct {
	mutex    sync.Mutex
	storage  map[string]string
	userData map[string][]UserData
	file     *os.File
	encoder  *json.Encoder
}

type storageData struct {
	ShortURL string `json:"shortUrl"`
	FullURL  string `json:"fullUrl"`
	UserID   string `json:"userId"`
}

func NewFileStorage(filepath string) (Storage, io.Closer, error) {
	file, openFileErr := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if openFileErr != nil {
		return nil, nil, openFileErr
	}

	storage := make(map[string]string)
	reader := bufio.NewReader(file)
	userData := make(map[string][]UserData)

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
		userData[data.UserID] = append(userData[data.UserID], UserData{
			FullURL:  data.FullURL,
			ShortURL: data.ShortURL,
		})
	}

	return &fileStorage{
		mutex:    sync.Mutex{},
		storage:  storage,
		userData: userData,
		file:     file,
		encoder:  json.NewEncoder(file),
	}, file, nil
}

func (s *fileStorage) SaveURL(ctx context.Context, input URLInput) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.storage[input.ShortURL] = input.FullURL
	s.userData[input.UserID] = append(s.userData[input.UserID], UserData{
		ShortURL: input.ShortURL,
		FullURL:  input.FullURL,
	})

	data := &storageData{
		ShortURL: input.ShortURL,
		FullURL:  input.FullURL,
		UserID:   input.UserID,
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

func (s *fileStorage) GetURLsByUserID(ctx context.Context, userID string) ([]UserData, error) {
	return s.userData[userID], nil
}

func (s *fileStorage) SaveBatch(ctx context.Context, batchInput []URLInput) error {
	return errors.New("method not implemented")
}

func (s *fileStorage) GetByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	return "", errors.New("method not implemented")
}

func (s *fileStorage) DeleteBatch(input []DeleteURLInput) error {
	return errors.New("method not implemented")
}
