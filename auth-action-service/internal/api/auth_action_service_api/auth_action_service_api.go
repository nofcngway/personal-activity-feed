package auth_action_service_api

import (
	"context"

	"github.com/nofcngway/auth-action-service/internal/pb/auth_action_api"
)

type authService interface {
	Register(ctx context.Context, username, password string) (token string, userID int64, err error)
	Login(ctx context.Context, username, password string) (token string, userID int64, err error)
	Logout(ctx context.Context, token string) error
	CreatePost(ctx context.Context, token string, postID int64) error
	Like(ctx context.Context, token string, postID int64) error
	Follow(ctx context.Context, token string, targetUserID int64) error
}

type AuthService = authService

// AuthActionServiceAPI реализует grpc AuthActionServiceServer (транспортный слой)
type AuthActionServiceAPI struct {
	auth_action_api.UnimplementedAuthActionServiceServer
	authService AuthService
}

func NewAuthActionServiceAPI(authService AuthService) *AuthActionServiceAPI {
	return &AuthActionServiceAPI{authService: authService}
}
