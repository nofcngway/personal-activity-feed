package auth_action_service_api

import (
	"context"

	"github.com/nofcngway/auth-action-service/internal/pb/models"
)

func (a *AuthActionServiceAPI) Logout(ctx context.Context, _ *models.Empty) (*models.ActionResponse, error) {
	token, err := tokenFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.authService.Logout(ctx, token); err != nil {
		return nil, mapErr(err)
	}
	return &models.ActionResponse{Ok: true}, nil
}
