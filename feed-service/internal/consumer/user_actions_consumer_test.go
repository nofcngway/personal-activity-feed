package consumer

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/segmentio/kafka-go"
)

type fakeReader struct {
	mu   sync.Mutex
	msgs []kafka.Message
	i    int
}

func (r *fakeReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.i >= len(r.msgs) {
		return kafka.Message{}, context.Canceled
	}
	m := r.msgs[r.i]
	r.i++
	return m, nil
}

func (r *fakeReader) Close() error { return nil }

type fakeReaderClose struct {
	closed bool
	err    error
}

func (r *fakeReaderClose) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return kafka.Message{}, context.Canceled
}

func (r *fakeReaderClose) Close() error {
	r.closed = true
	return r.err
}

type fakeReaderErr struct {
	err error
}

func (r *fakeReaderErr) ReadMessage(ctx context.Context) (kafka.Message, error) { return kafka.Message{}, r.err }
func (r *fakeReaderErr) Close() error                                           { return nil }

type fakeProcessor struct {
	mu     sync.Mutex
	calls  int
	values [][]byte
	errOn  map[int]error // call index -> error
}

func (p *fakeProcessor) Handle(ctx context.Context, payload []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	p.values = append(p.values, payload)
	if err := p.errOn[p.calls]; err != nil {
		return err
	}
	return nil
}

func TestConsume_CallsProcessor(t *testing.T) {
	t.Parallel()

	proc := &fakeProcessor{errOn: map[int]error{}}
	r := &fakeReader{
		msgs: []kafka.Message{{Value: []byte("m1")}},
	}
	c := &UserActionsConsumer{reader: r, processor: proc}

	err := c.Consume(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}

	if proc.calls != 1 {
		t.Fatalf("expected 1 call, got %d", proc.calls)
	}
	if string(proc.values[0]) != "m1" {
		t.Fatalf("unexpected payload: %q", string(proc.values[0]))
	}
}

func TestConsume_ProcessorErrorIsSwallowed(t *testing.T) {
	t.Parallel()

	proc := &fakeProcessor{errOn: map[int]error{1: errors.New("bad")}}
	r := &fakeReader{
		msgs: []kafka.Message{
			{Value: []byte("bad")},
			{Value: []byte("good")},
		},
	}
	c := &UserActionsConsumer{reader: r, processor: proc}

	err := c.Consume(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
	if proc.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", proc.calls)
	}
	if string(proc.values[1]) != "good" {
		t.Fatalf("expected second payload to be 'good', got %q", string(proc.values[1]))
	}
}

func TestClose_DelegatesToReader(t *testing.T) {
	t.Parallel()

	r := &fakeReaderClose{err: errors.New("close")}
	c := &UserActionsConsumer{reader: r, processor: &fakeProcessor{errOn: map[int]error{}}}
	err := c.Close()
	if err == nil || err.Error() != "close" {
		t.Fatalf("expected close error, got: %v", err)
	}
	if !r.closed {
		t.Fatalf("expected reader to be closed")
	}
}

func TestConsume_ReaderErrorIsReturned(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("dial failed")
	r := &fakeReaderErr{err: sentinel}
	proc := &fakeProcessor{errOn: map[int]error{}}
	c := &UserActionsConsumer{reader: r, processor: proc}

	err := c.Consume(context.Background())
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel err, got: %v", err)
	}
	if proc.calls != 0 {
		t.Fatalf("expected processor not called")
	}
}


