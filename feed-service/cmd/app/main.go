package main

import (
	"fmt"
	"os"

	"github.com/nofcngway/feed-service/config"
	"github.com/nofcngway/feed-service/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	storage := bootstrap.InitPGStorage(cfg)
	feedService := bootstrap.InitFeedService(storage)
	processor := bootstrap.InitUserActionsProcessor(feedService)
	api := bootstrap.InitFeedAPI(feedService)
	consumer := bootstrap.InitUserActionsConsumer(cfg, processor)

	bootstrap.AppRun(cfg, api, consumer, storage)
}
