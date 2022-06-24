package storage

import (
	"context"
	"errors"
)

var ErrAlreadyExist = errors.New("original url already exist")
var ErrIsDeleted = errors.New("is deleted")

type URLInput struct {
	ShortURL string
	FullURL  string
	UserID   string
}

type UserData struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

type DeleteURLInput struct {
	UserID string
	URL    string
}

type Storage interface {
	SaveURL(ctx context.Context, input URLInput) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	GetURLsByUserID(ctx context.Context, userID string) ([]UserData, error)
	SaveBatch(ctx context.Context, batchInput []URLInput) error
	GetByOriginalURL(ctx context.Context, originalURL string) (string, error)
	DeleteBatch(input []DeleteURLInput) error
}
