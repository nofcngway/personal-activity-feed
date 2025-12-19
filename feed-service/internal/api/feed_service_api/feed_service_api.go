package feed_service_api

import (
	"context"

	"github.com/nofcngway/feed-service/internal/pb/feed_api"
	"github.com/nofcngway/feed-service/internal/storage/pgstorage"
)

type feedService interface {
	GetFeed(ctx context.Context, userID int64, limit, offset int32) ([]pgstorage.FeedItem, error)
}

type FeedService = feedService

// FeedServiceAPI реализует grpc FeedServiceServer (транспортный слой)
type FeedServiceAPI struct {
	feed_api.UnimplementedFeedServiceServer
	feedService FeedService
}

func NewFeedServiceAPI(feedService FeedService) *FeedServiceAPI {
	return &FeedServiceAPI{feedService: feedService}
}
