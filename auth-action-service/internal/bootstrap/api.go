package bootstrap

import (
	server "github.com/nofcngway/auth-action-service/internal/api/auth_action_service_api"
	"github.com/nofcngway/auth-action-service/internal/services/authservice"
)

func InitAuthAPI(authService *authservice.Service) *server.AuthActionServiceAPI {
	return server.NewAuthActionServiceAPI(authService)
}
