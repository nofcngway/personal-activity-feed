package auth_action_service_api

import (
	"context"

	"github.com/nofcngway/auth-action-service/internal/pb/models"
)

func (a *AuthActionServiceAPI) CreatePost(ctx context.Context, req *models.CreatePostRequest) (*models.ActionResponse, error) {
	token, err := tokenFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.authService.CreatePost(ctx, token, req.GetPostId()); err != nil {
		return nil, mapErr(err)
	}
	return &models.ActionResponse{Ok: true}, nil
}
