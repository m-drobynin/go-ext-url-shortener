package model

import (
	"context"
	"database/sql"
	"uuid"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(db *sql.DB, connectionString string) (*PostgresDB, error) {
	res := &PostgresDB{
		db: db,
	}

	return res, nil
}

func (psDb *PostgresDB) TryPut(ctx context.Context, key string, value string) error {
	_, err := psDb.db.ExecContext(ctx, "INSERT INTO urls (uuid, short_url, original_url) VALUES ($1, $2, $3)", uuid.NewV4().String(), key, value)

	return err
}

func (psDb *PostgresDB) Get(ctx context.Context, key string) (*string, error) {
	row := psDb.db.QueryRowContext(ctx, "SELECT uuid, short_url, original_url FROM urls WHERE short_url=$1", key)

	var urlRow ShortenedURLRow

	err := row.Scan(&urlRow.UUID, &urlRow.ShortURL, &urlRow.OriginalURL)

	if err != nil {
		return nil, err
	}

	return &urlRow.OriginalURL, nil
}
