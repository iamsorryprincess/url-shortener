package service

import (
	"math/rand"
	"sync"
	"time"
)

type Storage interface {
	SaveURL(url string, shortURL string) error
	GetURL(shortURL string) string
}

type UserData struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

type URLService struct {
	storage      Storage
	randomizer   *rand.Rand
	randomMatrix []string
	userMutex    sync.Mutex
	userUrls     map[string][]UserData
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
		userUrls:  make(map[string][]UserData),
		userMutex: sync.Mutex{},
	}
}

func (service *URLService) SaveURL(url string, userID string, baseURL string) (string, error) {
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

	service.userMutex.Lock()
	service.userUrls[userID] = append(service.userUrls[userID], UserData{
		FullURL:  url,
		ShortURL: baseURL + "/" + key,
	})
	service.userMutex.Unlock()

	if err != nil {
		return "", err
	}

	return key, nil
}

func (service *URLService) GetURL(url string) string {
	return service.storage.GetURL(url)
}

func (service *URLService) GetUserData(userID string) []UserData {
	return service.userUrls[userID]
}
