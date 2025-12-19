package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	topic  string
	writer writer
}

type writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type UserActionEvent struct {
	EventID   string    `json:"event_id"`
	UserID    int64     `json:"user_id"`
	Action    string    `json:"action"`
	TargetID  int64     `json:"target_id"`
	Timestamp time.Time `json:"timestamp"`
}

func New(brokers []string, topic string) *Producer {
	return &Producer{
		topic: topic,
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) Close() error { return p.writer.Close() }

func (p *Producer) Publish(ctx context.Context, userID int64, action string, targetID int64) error {
	ev := UserActionEvent{
		EventID:   uuid.NewString(),
		UserID:    userID,
		Action:    action,
		TargetID:  targetID,
		Timestamp: time.Now().UTC(),
	}

	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(ev.EventID),
		Value: b,
		Time:  ev.Timestamp,
	}
	return p.writer.WriteMessages(ctx, msg)
}
