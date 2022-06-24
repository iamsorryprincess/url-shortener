package worker

import (
	"log"

	"github.com/iamsorryprincess/url-shortener/internal/storage"
)

type Worker struct {
	balancerChannel chan storage.DeleteURLInput
	inputChannels   []chan storage.DeleteURLInput
	storage         storage.Storage
}

func NewWorker(storage storage.Storage) *Worker {
	return &Worker{
		storage: storage,
	}
}

func (w *Worker) Start(workersCount int, poolSize int) {
	go func(workersCount int, poolSize int) {
		w.inputChannels = make([]chan storage.DeleteURLInput, workersCount)

		for i := 0; i < workersCount; i++ {
			w.inputChannels[i] = make(chan storage.DeleteURLInput)
			go func(poolSize int, ch chan storage.DeleteURLInput) {
				pool := make([]storage.DeleteURLInput, poolSize)
				count := 0
				for urlData := range ch {
					pool[count] = urlData
					if count == poolSize-1 {
						err := w.storage.DeleteBatch(pool)
						if err != nil {
							log.Println(err)
						}
						count = 0
						continue
					}
					count++
				}
			}(poolSize, w.inputChannels[i])
		}

		w.balancerChannel = make(chan storage.DeleteURLInput)
		count := 0

		for urlData := range w.balancerChannel {
			w.inputChannels[count] <- urlData
			if count == workersCount-1 {
				count = 0
				continue
			}
			count++
		}

		for _, ch := range w.inputChannels {
			close(ch)
		}
	}(workersCount, poolSize)
}

func (w *Worker) Process(userID string, urls []string) {
	go func(userID string, urls []string) {
		for _, url := range urls {
			w.balancerChannel <- storage.DeleteURLInput{
				UserID: userID,
				URL:    url,
			}
		}
	}(userID, urls)
}

func (w *Worker) Stop() {
	close(w.balancerChannel)
}
