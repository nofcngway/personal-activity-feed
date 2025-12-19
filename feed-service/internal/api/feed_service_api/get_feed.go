package feed_service_api

import (
	"context"
	"time"

	"github.com/nofcngway/feed-service/internal/pb/feed_api"
	"github.com/nofcngway/feed-service/internal/pb/models"
)

func (a *FeedServiceAPI) GetFeed(ctx context.Context, req *feed_api.GetFeedRequest) (*feed_api.GetFeedResponse, error) {
	items, err := a.feedService.GetFeed(ctx, req.GetUserId(), req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, err
	}

	resp := &feed_api.GetFeedResponse{
		Items: make([]*models.FeedItemModel, 0, len(items)),
	}
	for _, it := range items {
		resp.Items = append(resp.Items, &models.FeedItemModel{
			ActorId:   it.ActorID,
			Action:    it.Action,
			TargetId:  it.TargetID,
			CreatedAt: it.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	return resp, nil
}
