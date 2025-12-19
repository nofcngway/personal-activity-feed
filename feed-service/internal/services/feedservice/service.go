package feedservice

import (
	"context"
	"time"

	"github.com/nofcngway/feed-service/internal/storage/pgstorage"
)

type Storage interface {
	InsertFeedItem(ctx context.Context, userID, actorID int64, action string, targetID int64, createdAt time.Time) error
	GetFeed(ctx context.Context, userID int64, limit, offset int32) ([]pgstorage.FeedItem, error)
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) AddEvent(ctx context.Context, userID, actorID int64, action string, targetID int64, createdAt time.Time) error {
	return s.storage.InsertFeedItem(ctx, userID, actorID, action, targetID, createdAt)
}

func (s *Service) GetFeed(ctx context.Context, userID int64, limit, offset int32) ([]pgstorage.FeedItem, error) {
	return s.storage.GetFeed(ctx, userID, limit, offset)
}


