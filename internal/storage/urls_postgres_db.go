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

const insertURLSQL = `INSERT INTO urls (uuid, short_url, original_url, user_id) VALUES ($1, $2, $3, $4)`
const findURLSQL = `SELECT uuid, short_url, original_url FROM urls WHERE short_url=$1`
const findByOriginalURLSQL = `SELECT uuid, short_url, original_url FROM urls WHERE original_url=$1`
const findUrlsByUser = `SELECT uuid, short_url, original_url FROM urls WHERE user_id=$1`

type UrlsPostgresDB struct {
	db *sql.DB
}

func NewUrlsPostgresDB(db *sql.DB) (*UrlsPostgresDB, error) {
	res := &UrlsPostgresDB{
		db: db,
	}

	return res, nil
}

func (psDb *UrlsPostgresDB) Put(ctx context.Context, key string, value string, userID string) error {
	_, err := psDb.db.ExecContext(ctx, insertURLSQL, uuid.NewV4().String(), key, value, userID)

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

func (psDb *UrlsPostgresDB) PutBatch(ctx context.Context, pairs []model.ValuesPair, userID string) error {
	tx, err := psDb.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, insertURLSQL)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, pair := range pairs {
		_, err := stmt.ExecContext(ctx, uuid.NewV4().String(), pair.Key, pair.Value, userID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (psDb *UrlsPostgresDB) Get(ctx context.Context, key string) (*string, error) {
	row := psDb.db.QueryRowContext(ctx, findURLSQL, key)

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

func (psDb *UrlsPostgresDB) GetCode(ctx context.Context, originalURL string) (*string, error) {
	row := psDb.db.QueryRowContext(ctx, findByOriginalURLSQL, originalURL)

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

func (psDb *UrlsPostgresDB) GetAllByUserId(ctx context.Context, userID string) (*[]model.ValuesPair, error) {
	rows, err := psDb.db.QueryContext(ctx, findUrlsByUser, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrDBNotFound
		}

		return nil, err
	}

	defer rows.Close()

	var result []model.ValuesPair

	for rows.Next() {
		var urlRow model.ShortenedURLRow
		err = rows.Scan(&urlRow.UUID, &urlRow.ShortURL, &urlRow.OriginalURL)

		if err != nil {
			return nil, err
		}

		result = append(result, model.ValuesPair{
			Key:   urlRow.ShortURL,
			Value: urlRow.OriginalURL,
		})
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return &result, nil
}
