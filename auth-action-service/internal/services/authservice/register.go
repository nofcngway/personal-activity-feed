package authservice

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, username, password string) (token string, userID int64, err error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", 0, ErrInvalidArgument
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", 0, err
	}

	userID, err = s.storage.CreateUser(ctx, username, string(hash))
	if err != nil {
		// упрощенно: считаем, что любая ошибка создания = already exists/конфликт
		return "", 0, ErrUserAlreadyExists
	}

	token, err = s.createSession(ctx, userID)
	if err != nil {
		return "", 0, err
	}
	return token, userID, nil
}


