package migrations

import (
	"context"
	"database/sql"
	"os"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	file, err := os.OpenFile("internal/storage/migrations/initial.sql", os.O_RDONLY, 0777)

	if err != nil {
		return err
	}

	defer file.Close()
	info, err := file.Stat()

	if err != nil {
		return err
	}

	bytes := make([]byte, info.Size())
	_, err = file.Read(bytes)

	if err != nil {
		return err
	}

	sqlMigration := string(bytes)
	_, err = db.ExecContext(ctx, sqlMigration)
	return err
}
