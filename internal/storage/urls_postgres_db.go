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

const fields = `uuid, short_url, original_url, user_id, deleted`

const insertURLSQL = `INSERT INTO urls (` + fields + `) VALUES ($1, $2, $3, $4, $5)`
const findURLSQL = `SELECT ` + fields + ` FROM urls WHERE short_url=$1`
const findByOriginalURLSQL = `SELECT short_url FROM urls WHERE original_url=$1`
const findUrlsByUser = `SELECT ` + fields + ` FROM urls WHERE user_id=$1`
const markDeletedURLSQL = `UPDATE urls SET deleted=true WHERE short_url=$1 AND user_id=$2`

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
	_, err := psDb.db.ExecContext(ctx, insertURLSQL, uuid.NewV4().String(), key, value, userID, false)

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

		return err
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
		_, err := stmt.ExecContext(ctx, uuid.NewV4().String(), pair.Key, pair.Value, userID, false)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (psDb *UrlsPostgresDB) Get(ctx context.Context, key string) (*model.URLRow, error) {
	row := psDb.db.QueryRowContext(ctx, findURLSQL, key)

	var urlRow model.URLRow
	var userID sql.NullString

	err := row.Scan(&urlRow.UUID, &urlRow.ShortURL, &urlRow.OriginalURL, userID, &urlRow.Deleted)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrDBNotFound
		}

		return nil, err
	}

	if urlRow.Deleted {
		return nil, model.ErrDBDeleted
	}

	if userID.Valid {
		urlRow.UserID = userID.String
	}

	return &urlRow, nil
}

func (psDb *UrlsPostgresDB) GetCode(ctx context.Context, originalURL string) (*string, error) {
	row := psDb.db.QueryRowContext(ctx, findByOriginalURLSQL, originalURL)

	var shortURL string
	err := row.Scan(&shortURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrDBNotFound
		}

		return nil, err
	}

	return &shortURL, nil
}

func (psDb *UrlsPostgresDB) GetAllByUserID(ctx context.Context, userID string) (*[]model.URLRow, error) {
	rows, err := psDb.db.QueryContext(ctx, findUrlsByUser, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrDBNotFound
		}

		return nil, err
	}

	defer rows.Close()

	var result []model.URLRow

	for rows.Next() {
		var urlRow model.URLRow
		err = rows.Scan(&urlRow.UUID, &urlRow.ShortURL, &urlRow.OriginalURL, &urlRow.UserID, &urlRow.Deleted)

		if err != nil {
			return nil, err
		}

		result = append(result, urlRow)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (psDb *UrlsPostgresDB) DeleteBatch(ctx context.Context, codes []string, userID string) error {
	tx, err := psDb.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, markDeletedURLSQL)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, code := range codes {
		_, err := stmt.ExecContext(ctx, code, userID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
