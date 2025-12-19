package authservice

import (
	"context"
	"errors"
	"strings"

	"github.com/nofcngway/auth-action-service/internal/storage/pgstorage"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(ctx context.Context, username, password string) (token string, userID int64, err error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", 0, ErrInvalidArgument
	}

	u, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgstorage.ErrUserNotFound) {
			return "", 0, ErrInvalidCredentials
		}
		return "", 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", 0, ErrInvalidCredentials
	}

	token, err = s.createSession(ctx, u.ID)
	if err != nil {
		return "", 0, err
	}
	return token, u.ID, nil
}


