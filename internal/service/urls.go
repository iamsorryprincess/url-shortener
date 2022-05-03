package service

import (
	"math/rand"
	"time"
)

type Storage interface {
	SaveURL(url string, shortURL string) error
	GetURL(shortURL string) string
}

type URLService struct {
	storage      Storage
	randomizer   *rand.Rand
	randomMatrix []string
}

func NewURLService(storage Storage) *URLService {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return &URLService{
		storage:    storage,
		randomizer: r1,
		randomMatrix: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9",
			"A", "B", "C", "D", "E", "F", "G", "H", "L",
			"M", "N", "O", "P", "Q", "R", "T", "S", "U",
			"V", "W", "X", "Y", "Z"},
	}
}

func (service *URLService) SaveURL(url string) (string, error) {
	n := len(service.randomMatrix) - 1
	key := ""
	existingKey := "1"

	for existingKey != "" {
		key = ""
		for i := 1; i <= 10; i++ {
			index := service.randomizer.Intn(n)
			key += service.randomMatrix[index]
		}
		existingKey = service.storage.GetURL(key)
	}

	err := service.storage.SaveURL(url, key)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (service *URLService) GetURL(url string) string {
	return service.storage.GetURL(url)
}
