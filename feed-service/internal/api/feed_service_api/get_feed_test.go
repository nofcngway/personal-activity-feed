package feed_service_api

import (
	"context"
	"testing"
	"time"

	api_mocks "github.com/nofcngway/feed-service/internal/api/feed_service_api/mocks"
	"github.com/nofcngway/feed-service/internal/pb/feed_api"
	"github.com/nofcngway/feed-service/internal/storage/pgstorage"
	"github.com/stretchr/testify/mock"
)

func TestFeedServiceAPI_GetFeed_MapsItems(t *testing.T) {
	t.Parallel()

	ts := time.Date(2025, 12, 13, 10, 0, 0, 0, time.UTC)
	items := []pgstorage.FeedItem{
		{ActorID: 42, Action: "like", TargetID: 101, CreatedAt: ts},
	}

	svc := api_mocks.NewFeedService(t)
	svc.EXPECT().GetFeed(mock.Anything, int64(1), int32(10), int32(0)).Return(items, nil).Once()

	api := NewFeedServiceAPI(svc)
	resp, err := api.GetFeed(context.Background(), &feed_api.GetFeedRequest{
		UserId: 1,
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].ActorId != 42 || resp.Items[0].Action != "like" || resp.Items[0].TargetId != 101 || resp.Items[0].CreatedAt != ts.Format(time.RFC3339) {
		t.Fatalf("unexpected item: %+v", resp.Items[0])
	}
}
