package storage

import (
	"context"
	"errors"
)

var ErrAlreadyExist = errors.New("original url already exist")

type URLInput struct {
	ShortURL string
	FullURL  string
}

type Storage interface {
	SaveURL(ctx context.Context, url string, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	SaveBatch(ctx context.Context, batchInput []URLInput) error
	GetByOriginalURL(ctx context.Context, originalURL string) (string, error)
}
