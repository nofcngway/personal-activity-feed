package bootstrap

import (
	"github.com/nofcngway/feed-service/config"
	"github.com/nofcngway/feed-service/internal/consumer"
	useractionsprocessor "github.com/nofcngway/feed-service/internal/services/processors/user_actions_processor"
)

func InitUserActionsConsumer(cfg *config.Config, processor *useractionsprocessor.UserActionsProcessor) *consumer.UserActionsConsumer {
	return consumer.NewUserActionsConsumer(cfg.Kafka.Brokers, cfg.Kafka.TopicName, cfg.Kafka.GroupID, processor)
}


