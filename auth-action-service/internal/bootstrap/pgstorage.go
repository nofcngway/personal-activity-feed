package bootstrap

import (
	"fmt"
	"log"

	"github.com/nofcngway/auth-action-service/config"
	"github.com/nofcngway/auth-action-service/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGStorage {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	storage, err := pgstorage.NewPGStorage(connString)
	if err != nil {
		log.Panicf("ошибка инициализации БД, %v", err)
		panic(err)
	}
	return storage
}
