package worker

import "github.com/iamsorryprincess/url-shortener/internal/storage"

type URLInput struct {
	UserID string
	UrlID  string
}

type Worker struct {
	inputCh chan URLInput
	storage storage.Storage
}

func NewWorker(storage storage.Storage) *Worker {
	return &Worker{
		storage: storage,
	}
}

func (w *Worker) Push(input URLInput) {
	w.inputCh <- input
}

func (w *Worker) Start(workersCount int, poolSize int) {
	w.inputCh = make(chan URLInput)
	for i := 0; i < workersCount; i++ {
		go func() {
			pool := make([]URLInput, poolSize)

			for {
				select {
				case urlData, isOk := <-w.inputCh:
					if !isOk {
						w.inputCh = nil
						return
					}
					pool = append(pool, urlData)
					if len(pool) == poolSize {
						// delete batch
						pool = pool[:0]
					}
				}
			}
		}()
	}
}

func (w *Worker) Close() {
	close(w.inputCh)
}
