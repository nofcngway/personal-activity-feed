package pgstorage

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const usersTable = "users"

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

var ErrUserNotFound = errors.New("user not found")

func (s *PGStorage) CreateUser(ctx context.Context, username, passwordHash string) (int64, error) {
	q := squirrel.Insert(usersTable).
		Columns("username", "password_hash").
		Values(username, passwordHash).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := q.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	if err := s.db.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PGStorage) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	q := squirrel.Select("id", "username", "password_hash").
		From(usersTable).
		Where(squirrel.Eq{"username": username}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow(ctx, sql, args...)
	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}


