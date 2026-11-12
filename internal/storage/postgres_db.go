package storage

import (
	"context"
	"database/sql"
	"errors"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const insertSQL = `INSERT INTO urls (uuid, short_url, original_url) VALUES ($1, $2, $3)`
const findSQL = `SELECT uuid, short_url, original_url FROM urls WHERE short_url=$1`
const findByValueSQL = `SELECT uuid, short_url, original_url FROM urls WHERE original_url=$1`

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

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			switch pgErr.ConstraintName {
			case "idx_urls_original_url":
				return model.ErrUrlsConflict
			case "idx_urls_short_url":
				return model.ErrCodeConflict
			default:
				return err
			}
		}
	}

	return nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrDBNotFound
		}

		return nil, err
	}

	return &urlRow.OriginalURL, nil
}

func (psDb *PostgresDB) GetCode(ctx context.Context, originalURL string) (*string, error) {
	row := psDb.db.QueryRowContext(ctx, findByValueSQL, originalURL)

	var urlRow model.ShortenedURLRow

	err := row.Scan(&urlRow.UUID, &urlRow.ShortURL, &urlRow.OriginalURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrDBNotFound
		}

		return nil, err
	}

	return &urlRow.ShortURL, nil
}
