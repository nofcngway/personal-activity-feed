package user_actions_processor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nofcngway/feed-service/internal/services/feedservice"
)

type UserActionEvent struct {
	EventID   string    `json:"event_id"`
	UserID    int64     `json:"user_id"`
	Action    string    `json:"action"`
	TargetID  int64     `json:"target_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserActionsProcessor struct {
	feedService *feedservice.Service
}

func New(feedService *feedservice.Service) *UserActionsProcessor {
	return &UserActionsProcessor{feedService: feedService}
}

// Handle принимает сырой payload Kafka (json), валидирует/нормализует и пишет в ленту.
func (p *UserActionsProcessor) Handle(ctx context.Context, payload []byte) error {
	var ev UserActionEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		return err
	}

	createdAt := ev.Timestamp
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	// упрощение: лента "для себя"
	userID := ev.UserID
	actorID := ev.UserID

	return p.feedService.AddEvent(ctx, userID, actorID, ev.Action, ev.TargetID, createdAt)
}
