package user_actions_processor

import (
	"context"
	"testing"
	"time"

	"github.com/nofcngway/feed-service/internal/services/feedservice"
	"github.com/nofcngway/feed-service/internal/storage/pgstorage"
)

type storageMock struct {
	last struct {
		userID    int64
		actorID   int64
		action    string
		targetID  int64
		createdAt time.Time
	}
}

func (s *storageMock) InsertFeedItem(ctx context.Context, userID, actorID int64, action string, targetID int64, createdAt time.Time) error {
	s.last.userID = userID
	s.last.actorID = actorID
	s.last.action = action
	s.last.targetID = targetID
	s.last.createdAt = createdAt
	return nil
}

func (s *storageMock) GetFeed(ctx context.Context, userID int64, limit, offset int32) ([]pgstorage.FeedItem, error) {
	return nil, nil
}

func TestProcessor_Handle_ParsesAndWrites(t *testing.T) {
	t.Parallel()

	st := &storageMock{}
	svc := feedservice.New(st)
	p := New(svc)

	ts := time.Date(2025, 12, 13, 10, 0, 0, 0, time.UTC)
	payload := []byte(`{"event_id":"e1","user_id":42,"action":"like","target_id":101,"timestamp":"2025-12-13T10:00:00Z"}`)

	if err := p.Handle(context.Background(), payload); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if st.last.userID != 42 || st.last.actorID != 42 || st.last.action != "like" || st.last.targetID != 101 || !st.last.createdAt.Equal(ts) {
		t.Fatalf("unexpected stored: %+v", st.last)
	}
}

func TestProcessor_Handle_SetsTimestampIfMissing(t *testing.T) {
	t.Parallel()

	st := &storageMock{}
	svc := feedservice.New(st)
	p := New(svc)

	before := time.Now().UTC()
	payload := []byte(`{"event_id":"e1","user_id":1,"action":"create_post","target_id":999}`)

	if err := p.Handle(context.Background(), payload); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if st.last.createdAt.IsZero() {
		t.Fatalf("expected createdAt to be set")
	}
	if st.last.createdAt.Before(before.Add(-2*time.Second)) || st.last.createdAt.After(time.Now().UTC().Add(2*time.Second)) {
		t.Fatalf("createdAt not within expected range: %s", st.last.createdAt)
	}
	if st.last.createdAt.Location() != time.UTC {
		t.Fatalf("expected UTC location, got: %v", st.last.createdAt.Location())
	}
}

func TestProcessor_Handle_InvalidJSON_ReturnsError(t *testing.T) {
	t.Parallel()

	st := &storageMock{}
	svc := feedservice.New(st)
	p := New(svc)

	if err := p.Handle(context.Background(), []byte(`{not-json`)); err == nil {
		t.Fatalf("expected error")
	}
}
