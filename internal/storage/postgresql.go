package storage

import (
	"context"
	"database/sql"

	"github.com/iamsorryprincess/url-shortener/internal/storage/migrations"
)

type postgresqlStorage struct {
	db *sql.DB
}

func NewPostgresqlStorage(db *sql.DB) (Storage, error) {
	err := migrations.Migrate(context.Background(), db)

	if err != nil {
		return nil, err
	}

	return &postgresqlStorage{
		db: db,
	}, nil
}

func (s *postgresqlStorage) SaveURL(ctx context.Context, url string, shortURL string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO public.urls (short_url, original_url)\nVALUES ($1, $2);", shortURL, url)

	if err != nil {
		return err
	}

	return nil
}

func (s *postgresqlStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	result := ""
	err := s.db.QueryRowContext(ctx, "SELECT original_url FROM public.urls WHERE short_url=$1", shortURL).Scan(&result)

	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *postgresqlStorage) SaveBatch(ctx context.Context, input []URLInput) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO public.urls (short_url, original_url) VALUES ($1, $2)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, inputData := range input {
		_, err = stmt.ExecContext(ctx, inputData.ShortUrl, inputData.FullUrl)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
