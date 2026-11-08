package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrDBConnStringEmpty = errors.New("db_conn_string_empty")

func PreprareDBConnection(connString string) (*sql.DB, error) {
	if connString == "" {
		return nil, ErrDBConnStringEmpty
	}

	db, err := sql.Open("pgx", connString)

	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "pgx", driver)

	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	return db, nil
}
