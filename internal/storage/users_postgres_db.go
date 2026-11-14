package storage

import (
	"context"
	"database/sql"
	"errors"
	"m-drobynin/go-ext-url-shortener/internal/model"
)

const insertUserSQL = `INSERT INTO users (user_id) VALUES ($1)`
const findUserSQL = `SELECT user_id FROM users WHERE user_id=$1`

type UsersPostgresDB struct {
	db      *sql.DB
	usersDb *UsersMemoryDatabase
}

func NewUsersPostgresDB(db *sql.DB) (*UsersPostgresDB, error) {
	res := &UsersPostgresDB{
		db: db,
	}

	return res, nil
}

func (psDb *UsersPostgresDB) Create(ctx context.Context) (*model.User, error) {
	user, err := psDb.usersDb.Create(ctx)

	if err != nil {
		return nil, err
	}

	_, err = psDb.db.ExecContext(ctx, insertUserSQL, user.UserID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (psDb *UsersPostgresDB) GetByID(ctx context.Context, userID string) (*model.User, error) {
	row := psDb.db.QueryRowContext(ctx, findUserSQL, userID)

	var user model.User

	err := row.Scan(&user.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrDBNotFound
		}

		return nil, err
	}

	return &user, nil
}
