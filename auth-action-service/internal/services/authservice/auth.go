package authservice

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

func (s *Service) authorize(ctx context.Context, token string) (int64, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0, ErrUnauthorized
	}
	sess, err := s.sessions.Get(ctx, token)
	if err != nil {
		if err == redis.Nil {
			return 0, ErrUnauthorized
		}
		return 0, err
	}
	return sess.UserID, nil
}


