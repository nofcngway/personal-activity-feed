package producer

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

type fakeWriter struct {
	msgs []kafka.Message
	err  error
}

func (w *fakeWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	w.msgs = append(w.msgs, msgs...)
	return w.err
}

func (w *fakeWriter) Close() error { return nil }

func TestProducer_Publish_WritesMessage(t *testing.T) {
	t.Parallel()

	w := &fakeWriter{}
	p := &Producer{topic: "user-actions", writer: w}

	if err := p.Publish(context.Background(), 42, "like", 101); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(w.msgs) != 1 {
		t.Fatalf("expected 1 msg, got %d", len(w.msgs))
	}

	var ev UserActionEvent
	if err := json.Unmarshal(w.msgs[0].Value, &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.UserID != 42 || ev.Action != "like" || ev.TargetID != 101 {
		t.Fatalf("unexpected event: %+v", ev)
	}
	if ev.EventID == "" {
		t.Fatalf("expected event_id")
	}
	if string(w.msgs[0].Key) != ev.EventID {
		t.Fatalf("expected key==event_id, key=%q event_id=%q", string(w.msgs[0].Key), ev.EventID)
	}
	if ev.Timestamp.IsZero() || ev.Timestamp.Location() != time.UTC {
		t.Fatalf("expected UTC timestamp, got %v", ev.Timestamp)
	}
	if !w.msgs[0].Time.Equal(ev.Timestamp) {
		t.Fatalf("expected message time == event timestamp")
	}
}

func TestProducer_Publish_PropagatesWriterError(t *testing.T) {
	t.Parallel()

	want := errors.New("write failed")
	w := &fakeWriter{err: want}
	p := &Producer{topic: "user-actions", writer: w}

	err := p.Publish(context.Background(), 1, "like", 2)
	if !errors.Is(err, want) {
		t.Fatalf("expected writer error, got: %v", err)
	}
}

func TestProducer_Close_Delegates(t *testing.T) {
	t.Parallel()

	w := &closeWriter{}

	p := &Producer{topic: "t", writer: w}
	_ = p.Close()
	if !w.closed {
		t.Fatalf("expected close to be called")
	}
}

type closeWriter struct {
	closed   bool
	closeErr error
}

func (w *closeWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error { return nil }
func (w *closeWriter) Close() error {
	w.closed = true
	return w.closeErr
}
