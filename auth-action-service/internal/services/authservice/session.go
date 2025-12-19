package authservice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nofcngway/auth-action-service/internal/sessions"
)

func (s *Service) createSession(ctx context.Context, userID int64) (string, error) {
	token := uuid.NewString()
	err := s.sessions.Set(ctx, token, sessions.Session{
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}, s.sessionTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}


