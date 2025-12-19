package auth_action_service_api

import (
	"context"

	"github.com/nofcngway/auth-action-service/internal/pb/models"
)

func (a *AuthActionServiceAPI) Follow(ctx context.Context, req *models.FollowRequest) (*models.ActionResponse, error) {
	token, err := tokenFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.authService.Follow(ctx, token, req.GetUserId()); err != nil {
		return nil, mapErr(err)
	}
	return &models.ActionResponse{Ok: true}, nil
}
