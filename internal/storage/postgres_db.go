package storage

import (
	"context"
	"database/sql"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"uuid"
)

const insertSQL = `INSERT INTO urls (uuid, short_url, original_url) VALUES ($1, $2, $3)`
const findSQL = `SELECT uuid, short_url, original_url FROM urls WHERE short_url=$1`

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(db *sql.DB, connectionString string) (*PostgresDB, error) {
	res := &PostgresDB{
		db: db,
	}

	return res, nil
}

func (psDb *PostgresDB) Put(ctx context.Context, key string, value string) error {
	_, err := psDb.db.ExecContext(ctx, insertSQL, uuid.NewV4().String(), key, value)

	return err
}

func (psDb *PostgresDB) PutBatch(ctx context.Context, pairs []model.ValuesPair) error {
	tx, err := psDb.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, insertSQL)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, pair := range pairs {
		_, err := stmt.ExecContext(ctx, uuid.NewV4().String(), pair.Key, pair.Value)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (psDb *PostgresDB) Get(ctx context.Context, key string) (*string, error) {
	row := psDb.db.QueryRowContext(ctx, findSQL, key)

	var urlRow model.ShortenedURLRow

	err := row.Scan(&urlRow.UUID, &urlRow.ShortURL, &urlRow.OriginalURL)

	if err != nil {
		return nil, err
	}

	return &urlRow.OriginalURL, nil
}
