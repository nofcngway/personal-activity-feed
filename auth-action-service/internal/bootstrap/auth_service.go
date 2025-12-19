package bootstrap

import (
	"time"

	"github.com/nofcngway/auth-action-service/config"
	kafkaproducer "github.com/nofcngway/auth-action-service/internal/kafka/producer"
	"github.com/nofcngway/auth-action-service/internal/services/authservice"
	"github.com/nofcngway/auth-action-service/internal/sessions"
	"github.com/nofcngway/auth-action-service/internal/storage/pgstorage"
)

func InitAuthService(cfg *config.Config, storage *pgstorage.PGStorage, sess *sessions.RedisStore, producer *kafkaproducer.Producer) *authservice.Service {
	ttl := time.Duration(cfg.Security.SessionTTLSeconds) * time.Second
	return authservice.New(storage, sess, producer, ttl)
}
