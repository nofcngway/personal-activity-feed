package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGStorage struct {
	db *pgxpool.Pool
}

func NewPGStorage(connString string) (*PGStorage, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига БД")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения к БД")
	}

	st := &PGStorage{db: db}
	if err := st.initTables(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return st, nil
}

func (s *PGStorage) Close() { s.db.Close() }

func (s *PGStorage) initTables(ctx context.Context) error {
	sql := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);`, usersTable)
	_, err := s.db.Exec(ctx, sql)
	if err != nil {
		return errors.Wrap(err, "init users table")
	}
	return nil
}


