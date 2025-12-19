package bootstrap

import (
	"github.com/nofcngway/auth-action-service/config"
	kafkaproducer "github.com/nofcngway/auth-action-service/internal/kafka/producer"
)

// В auth-action-service нет Kafka consumer'ов (сервис только публикует события),
// но по аналогии со students (где wiring Kafka лежит в consumers.go) держим Kafka init здесь.
func InitKafkaProducer(cfg *config.Config) *kafkaproducer.Producer {
	return kafkaproducer.New(cfg.Kafka.Brokers, cfg.Kafka.TopicName)
}
