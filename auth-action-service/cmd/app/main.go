package main

import (
	"fmt"
	"os"

	"github.com/nofcngway/auth-action-service/config"
	"github.com/nofcngway/auth-action-service/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	storage := bootstrap.InitPGStorage(cfg)
	sessions := bootstrap.InitSessions(cfg)
	producer := bootstrap.InitKafkaProducer(cfg)
	authService := bootstrap.InitAuthService(cfg, storage, sessions, producer)
	api := bootstrap.InitAuthAPI(authService)

	bootstrap.AppRun(cfg, api, storage, sessions, producer)
}
