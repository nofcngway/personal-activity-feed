package authservice

import (
	"context"
	"errors"
	"testing"
	"time"

	svc_mocks "github.com/nofcngway/auth-action-service/internal/services/authservice/mocks"
	"github.com/nofcngway/auth-action-service/internal/sessions"
	"github.com/nofcngway/auth-action-service/internal/storage/pgstorage"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

func TestRegister_OK_CreatesUserAndSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	st := svc_mocks.NewUserStorage(t)
	ss := svc_mocks.NewSessionStore(t)
	prod := svc_mocks.NewProducer(t)

	st.EXPECT().
		CreateUser(mock.Anything, "u1", mock.Anything).
		Run(func(_ context.Context, _ string, passwordHash string) {
			if passwordHash == "" || passwordHash == "p1" {
				t.Fatalf("expected hashed password")
			}
		}).
		Return(int64(10), nil).
		Once()

	ss.EXPECT().
		Get(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, _ string) (*sessions.Session, error) { return nil, redis.Nil }).
		Maybe()

	var setToken string
	ss.EXPECT().
		Set(mock.Anything, mock.Anything, mock.Anything, 2*time.Hour).
		Run(func(_ context.Context, token string, sess sessions.Session, ttl time.Duration) {
			setToken = token
			if sess.UserID != 10 {
				t.Fatalf("expected session userID 10, got %d", sess.UserID)
			}
			if ttl != 2*time.Hour {
				t.Fatalf("expected ttl 2h, got %v", ttl)
			}
		}).
		Return(nil).
		Once()

	svc := New(st, ss, prod, 2*time.Hour)
	token, userID, err := svc.Register(ctx, "u1", "p1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if userID != 10 || token == "" {
		t.Fatalf("unexpected token/userID: token=%q userID=%d", token, userID)
	}
	if token != setToken {
		t.Fatalf("expected Set called with returned token=%q, got %q", token, setToken)
	}
}

func TestRegister_InvalidArgs(t *testing.T) {
	t.Parallel()

	svc := New(
		svc_mocks.NewUserStorage(t),
		svc_mocks.NewSessionStore(t),
		svc_mocks.NewProducer(t),
		time.Hour,
	)

	_, _, err := svc.Register(context.Background(), " ", "")
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	t.Parallel()

	st := svc_mocks.NewUserStorage(t)
	ss := svc_mocks.NewSessionStore(t)
	prod := svc_mocks.NewProducer(t)

	st.EXPECT().GetUserByUsername(mock.Anything, "u").Return((*pgstorage.User)(nil), pgstorage.ErrUserNotFound).Once()

	svc := New(st, ss, prod, time.Hour)

	_, _, err := svc.Login(context.Background(), "u", "p")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogout_EmptyToken_Unauthorized(t *testing.T) {
	t.Parallel()

	svc := New(
		svc_mocks.NewUserStorage(t),
		svc_mocks.NewSessionStore(t),
		svc_mocks.NewProducer(t),
		time.Hour,
	)

	if err := svc.Logout(context.Background(), " "); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestActions_PublishCalled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Now().UTC()
	st := svc_mocks.NewUserStorage(t)
	ss := svc_mocks.NewSessionStore(t)
	prod := svc_mocks.NewProducer(t)

	ss.EXPECT().Get(mock.Anything, "t1").Return(&sessions.Session{UserID: 7, CreatedAt: now}, nil).Maybe()

	prod.EXPECT().Publish(mock.Anything, int64(7), "create_post", int64(101)).Return(nil).Once()
	prod.EXPECT().Publish(mock.Anything, int64(7), "like", int64(101)).Return(nil).Once()
	prod.EXPECT().Publish(mock.Anything, int64(7), "follow", int64(55)).Return(nil).Once()

	svc := New(st, ss, prod, time.Hour)

	if err := svc.CreatePost(ctx, "t1", 101); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := svc.Like(ctx, "t1", 101); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := svc.Follow(ctx, "t1", 55); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestAuthorize_RedisNil_Unauthorized(t *testing.T) {
	t.Parallel()

	ss := svc_mocks.NewSessionStore(t)
	ss.EXPECT().Get(mock.Anything, "t1").Return((*sessions.Session)(nil), redis.Nil).Once()

	svc := New(
		svc_mocks.NewUserStorage(t),
		ss,
		svc_mocks.NewProducer(t),
		time.Hour,
	)

	if err := svc.Like(context.Background(), "t1", 1); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestActions_InvalidIDs_ReturnInvalidArgument(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	ss := svc_mocks.NewSessionStore(t)
	ss.EXPECT().
		Get(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, _ string) (*sessions.Session, error) {
			return &sessions.Session{UserID: 1, CreatedAt: now}, nil
		}).
		Maybe()

	svc := New(
		svc_mocks.NewUserStorage(t),
		ss,
		svc_mocks.NewProducer(t),
		time.Hour,
	)

	if err := svc.CreatePost(context.Background(), "t", 0); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
	if err := svc.Like(context.Background(), "t", -1); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
	if err := svc.Follow(context.Background(), "t", 0); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

func TestActions_ProducerError_Propagates(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	ss := svc_mocks.NewSessionStore(t)
	ss.EXPECT().Get(mock.Anything, mock.Anything).Return(&sessions.Session{UserID: 1, CreatedAt: now}, nil).Maybe()
	want := errors.New("kafka down")
	prod := svc_mocks.NewProducer(t)
	prod.EXPECT().Publish(mock.Anything, int64(1), "like", int64(1)).Return(want).Once()

	svc := New(
		svc_mocks.NewUserStorage(t),
		ss,
		prod,
		time.Hour,
	)

	if err := svc.Like(context.Background(), "t", 1); !errors.Is(err, want) {
		t.Fatalf("expected producer error, got %v", err)
	}
}

func TestNew_DefaultSessionTTL_IsHour(t *testing.T) {
	t.Parallel()

	st := svc_mocks.NewUserStorage(t)
	ss := svc_mocks.NewSessionStore(t)
	prod := svc_mocks.NewProducer(t)

	st.EXPECT().CreateUser(mock.Anything, "u", mock.Anything).Return(int64(1), nil).Once()
	ss.EXPECT().
		Get(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, _ string) (*sessions.Session, error) { return nil, redis.Nil }).
		Maybe()

	gotTTL := time.Duration(0)
	ss.EXPECT().
		Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ string, _ sessions.Session, ttl time.Duration) { gotTTL = ttl }).
		Return(nil).
		Once()

	svc := New(st, ss, prod, 0)
	_, _, _ = svc.Register(context.Background(), "u", "p")
	if gotTTL != time.Hour {
		t.Fatalf("expected default ttl 1h, got %v", gotTTL)
	}
}
