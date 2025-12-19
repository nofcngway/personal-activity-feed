package auth_action_service_api

import (
	"errors"

	"github.com/nofcngway/auth-action-service/internal/services/authservice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapErr(err error) error {
	switch {
	case errors.Is(err, authservice.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, "invalid argument")
	case errors.Is(err, authservice.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "user already exists")
	case errors.Is(err, authservice.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, authservice.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, "unauthorized")
	default:
		return status.Error(codes.Internal, "internal")
	}
}
