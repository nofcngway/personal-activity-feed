package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type UserActionsConsumer struct {
	reader    reader
	processor processor
}

type reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type processor interface {
	Handle(ctx context.Context, payload []byte) error
}

func NewUserActionsConsumer(brokers []string, topic, groupID string, processor processor) *UserActionsConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &UserActionsConsumer{reader: r, processor: processor}
}

func (c *UserActionsConsumer) Close() error { return c.reader.Close() }

func (c *UserActionsConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if err := c.processor.Handle(ctx, msg.Value); err != nil {
			// упрощенно: пропускаем битые сообщения/ошибки обработки
			continue
		}
	}
}
