package bootstrap

import (
	"fmt"
	"log"

	"github.com/nofcngway/feed-service/config"
	"github.com/nofcngway/feed-service/internal/storage/pgstorage"
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

	st, err := pgstorage.NewPGStorage(connString)
	if err != nil {
		log.Panicf("ошибка инициализации БД, %v", err)
	}
	return st
}


