package feedservice

import (
	"context"
	"errors"
	"testing"
	"time"

	svc_mocks "github.com/nofcngway/feed-service/internal/services/feedservice/mocks"
	"github.com/nofcngway/feed-service/internal/storage/pgstorage"
	"github.com/stretchr/testify/mock"
)

func TestService_AddEvent_DelegatesToStorage(t *testing.T) {
	t.Parallel()

	st := svc_mocks.NewStorage(t)

	svc := New(st)
	now := time.Date(2025, 12, 13, 10, 0, 0, 0, time.UTC)
	st.EXPECT().InsertFeedItem(mock.Anything, int64(1), int64(2), "like", int64(123), now).Return(nil).Once()

	if err := svc.AddEvent(context.Background(), 1, 2, "like", 123, now); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestService_GetFeed_DelegatesToStorage(t *testing.T) {
	t.Parallel()

	want := []pgstorage.FeedItem{
		{ActorID: 1, Action: "create_post", TargetID: 101, CreatedAt: time.Now().UTC()},
	}

	st := svc_mocks.NewStorage(t)
	st.EXPECT().GetFeed(mock.Anything, int64(42), int32(20), int32(5)).Return(want, nil).Once()

	svc := New(st)
	got, err := svc.GetFeed(context.Background(), 42, 20, 5)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != len(want) || got[0].Action != want[0].Action {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestService_GetFeed_PropagatesError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("boom")
	st := svc_mocks.NewStorage(t)
	st.EXPECT().GetFeed(mock.Anything, int64(1), int32(10), int32(0)).Return(([]pgstorage.FeedItem)(nil), sentinel).Once()

	svc := New(st)
	_, err := svc.GetFeed(context.Background(), 1, 10, 0)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel err, got: %v", err)
	}
}
