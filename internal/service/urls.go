package service

import (
	"math/rand"
	"time"
)

type Storage interface {
	SaveURL(url string, shortURL string)
	GetURL(shortURL string) string
}

type URLService struct {
	storage      *Storage
	randomizer   *rand.Rand
	randomMatrix []string
}

func (service *URLService) SaveURL(url string) string {
	var n = len(service.randomMatrix) - 1
	var key = ""
	var existingKey = "1"

	for existingKey != "" {
		key = ""
		for i := 1; i <= 10; i++ {
			var index = service.randomizer.Intn(n)
			key += service.randomMatrix[index]
		}
		existingKey = (*service.storage).GetURL(url)
	}

	(*service.storage).SaveURL(url, key)
	return key
}

func (service *URLService) GetURL(url string) string {
	return (*service.storage).GetURL(url)
}

func InitURLService(storage *Storage) *URLService {
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
