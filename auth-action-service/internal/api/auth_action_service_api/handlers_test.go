package auth_action_service_api

import (
	"context"
	"errors"
	"testing"

	api_mocks "github.com/nofcngway/auth-action-service/internal/api/auth_action_service_api/mocks"
	"github.com/nofcngway/auth-action-service/internal/pb/models"
	"github.com/nofcngway/auth-action-service/internal/services/authservice"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ctxWithAuth(token string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", token))
}

func TestRegister_HandlerCallsService(t *testing.T) {
	t.Parallel()

	svc := api_mocks.NewAuthService(t)
	svc.EXPECT().Register(mock.Anything, "u", "p").Return("t1", int64(7), nil)

	api := NewAuthActionServiceAPI(svc)
	resp, err := api.Register(context.Background(), &models.AuthRequest{Username: "u", Password: "p"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.Token != "t1" || resp.UserId != 7 {
		t.Fatalf("unexpected resp: %+v", resp)
	}
}

func TestLogin_HandlerCallsService(t *testing.T) {
	t.Parallel()

	svc := api_mocks.NewAuthService(t)
	svc.EXPECT().Login(mock.Anything, "u2", "p2").Return("t2", int64(9), nil)

	api := NewAuthActionServiceAPI(svc)
	resp, err := api.Login(context.Background(), &models.AuthRequest{Username: "u2", Password: "p2"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.Token != "t2" || resp.UserId != 9 {
		t.Fatalf("unexpected resp: %+v", resp)
	}
}

func TestLogout_RequiresAuthHeader(t *testing.T) {
	t.Parallel()

	api := NewAuthActionServiceAPI(api_mocks.NewAuthService(t))
	_, err := api.Logout(context.Background(), &models.Empty{})
	st, _ := status.FromError(err)
	if st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got %v", st.Code())
	}
}

func TestLogout_PassesTokenToService(t *testing.T) {
	t.Parallel()

	svc := api_mocks.NewAuthService(t)
	svc.EXPECT().Logout(mock.Anything, "abc").Return(nil)
	api := NewAuthActionServiceAPI(svc)

	resp, err := api.Logout(ctxWithAuth("abc"), &models.Empty{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp == nil || !resp.Ok {
		t.Fatalf("expected ok response, got: %+v", resp)
	}
}

func TestLike_MapsServiceError(t *testing.T) {
	t.Parallel()

	svc := api_mocks.NewAuthService(t)
	svc.EXPECT().Like(mock.Anything, "t1", int64(1)).Return(authservice.ErrUnauthorized)
	api := NewAuthActionServiceAPI(svc)

	_, err := api.Like(ctxWithAuth("Bearer t1"), &models.LikeRequest{PostId: 1})
	st, _ := status.FromError(err)
	if st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got %v", st.Code())
	}
}

func TestLike_PassesTokenAndPostID(t *testing.T) {
	t.Parallel()

	svc := api_mocks.NewAuthService(t)
	svc.EXPECT().Like(mock.Anything, "t1", int64(123)).Return(nil)
	api := NewAuthActionServiceAPI(svc)

	resp, err := api.Like(ctxWithAuth("Bearer t1"), &models.LikeRequest{PostId: 123})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp == nil || !resp.Ok {
		t.Fatalf("expected ok response, got: %+v", resp)
	}
}

func TestFollow_PassesTokenAndUserID(t *testing.T) {
	t.Parallel()

	svc := api_mocks.NewAuthService(t)
	svc.EXPECT().Follow(mock.Anything, "t1", int64(77)).Return(nil)
	api := NewAuthActionServiceAPI(svc)
	resp, err := api.Follow(ctxWithAuth("Bearer t1"), &models.FollowRequest{UserId: 77})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp == nil || !resp.Ok {
		t.Fatalf("expected ok response, got: %+v", resp)
	}
}

func TestCreatePost_PassesTokenAndPostID(t *testing.T) {
	t.Parallel()

	svc := api_mocks.NewAuthService(t)
	svc.EXPECT().CreatePost(mock.Anything, "t1", int64(101)).Return(nil)
	api := NewAuthActionServiceAPI(svc)
	resp, err := api.CreatePost(ctxWithAuth("t1"), &models.CreatePostRequest{PostId: 101})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp == nil || !resp.Ok {
		t.Fatalf("expected ok response, got: %+v", resp)
	}
}

func TestRegister_MapsUnknownErrorToInternal(t *testing.T) {
	t.Parallel()

	svc := api_mocks.NewAuthService(t)
	svc.EXPECT().Register(mock.Anything, "u", "p").Return("", int64(0), errors.New("boom"))
	api := NewAuthActionServiceAPI(svc)

	_, err := api.Register(context.Background(), &models.AuthRequest{Username: "u", Password: "p"})
	st, _ := status.FromError(err)
	if st.Code() != codes.Internal {
		t.Fatalf("expected Internal, got %v", st.Code())
	}
}
