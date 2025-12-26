package authservice

import (
	"context"
	"strings"
)

func (s *Service) Logout(ctx context.Context, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return ErrUnauthorized
	}
	_ = s.sessions.Del(ctx, token)
	return nil
}

func (s *Service) CreatePost(ctx context.Context, token string, postID int64) error {
	userID, err := s.authorize(ctx, token)
	if err != nil {
		return err
	}
	if postID <= 0 {
		return ErrInvalidArgument
	}
	return s.producer.Publish(ctx, userID, "create_post", postID)
}

func (s *Service) Like(ctx context.Context, token string, postID int64) error {
	userID, err := s.authorize(ctx, token)
	if err != nil {
		return err
	}
	if postID <= 0 {
		return ErrInvalidArgument
	}
	return s.producer.Publish(ctx, userID, "like", postID)
}

func (s *Service) Follow(ctx context.Context, token string, targetUserID int64) error {
	userID, err := s.authorize(ctx, token)
	if err != nil {
		return err
	}
	if targetUserID <= 0 {
		return ErrInvalidArgument
	}

	exists, err := s.storage.UserExists(ctx, targetUserID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	return s.producer.Publish(ctx, userID, "follow", targetUserID)
}


