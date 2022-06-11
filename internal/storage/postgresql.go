package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iamsorryprincess/url-shortener/internal/storage/migrations"
	"github.com/jackc/pgx"
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

func (s *postgresqlStorage) SaveURL(ctx context.Context, input URLInput) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO public.urls (short_url, original_url, user_id) VALUES ($1, $2, $3);", input.ShortURL, input.FullURL, input.UserID)

	var pgError pgx.PgError

	if err != nil {
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return ErrAlreadyExist
		}
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
		_, err = stmt.ExecContext(ctx, inputData.ShortURL, inputData.FullURL)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *postgresqlStorage) GetByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	result := ""
	err := s.db.QueryRowContext(ctx, "SELECT short_url FROM public.urls WHERE original_url=$1", originalURL).Scan(&result)

	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *postgresqlStorage) GetURLsByUserID(ctx context.Context, userID string) ([]UserData, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT short_url, original_url FROM public.urls WHERE user_id=$1", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []UserData

	for rows.Next() {
		var userData UserData
		err = rows.Scan(&userData.ShortURL, &userData.FullURL)
		if err != nil {
			return nil, err
		}
		result = append(result, userData)
	}

	return result, nil
}
