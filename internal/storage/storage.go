package storage

import "context"

type URLInput struct {
	ShortUrl string
	FullUrl  string
}

type Storage interface {
	SaveURL(ctx context.Context, url string, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	SaveBatch(ctx context.Context, batchInput []URLInput) error
}
