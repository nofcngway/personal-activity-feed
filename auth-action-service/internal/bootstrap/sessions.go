package bootstrap

import (
	"github.com/nofcngway/auth-action-service/config"
	"github.com/nofcngway/auth-action-service/internal/sessions"
)

func InitSessions(cfg *config.Config) *sessions.RedisStore {
	return sessions.NewRedisStore(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
}


