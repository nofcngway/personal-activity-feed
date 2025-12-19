package authservice

import (
	"context"
	"time"

	"github.com/nofcngway/auth-action-service/internal/sessions"
	"github.com/nofcngway/auth-action-service/internal/storage/pgstorage"
)

type UserStorage interface {
	CreateUser(ctx context.Context, username, passwordHash string) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*pgstorage.User, error)
}

type SessionStore interface {
	Set(ctx context.Context, token string, session sessions.Session, ttl time.Duration) error
	Get(ctx context.Context, token string) (*sessions.Session, error)
	Del(ctx context.Context, token string) error
}

type Producer interface {
	Publish(ctx context.Context, userID int64, action string, targetID int64) error
	Close() error
}


