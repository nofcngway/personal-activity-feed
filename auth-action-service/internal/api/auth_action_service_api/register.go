package auth_action_service_api

import (
	"context"

	"github.com/nofcngway/auth-action-service/internal/pb/models"
)

func (a *AuthActionServiceAPI) Register(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error) {
	token, userID, err := a.authService.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, mapErr(err)
	}
	return &models.AuthResponse{Token: token, UserId: userID}, nil
}
