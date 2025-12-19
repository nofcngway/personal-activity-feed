package auth_action_service_api

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestTokenFromMetadata_Bearer(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer abc"))
	token, err := tokenFromMetadata(ctx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if token != "abc" {
		t.Fatalf("expected token 'abc', got %q", token)
	}
}

func TestTokenFromMetadata_RawToken(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "abc"))
	token, err := tokenFromMetadata(ctx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if token != "abc" {
		t.Fatalf("expected token 'abc', got %q", token)
	}
}

func TestTokenFromMetadata_Missing(t *testing.T) {
	t.Parallel()

	_, err := tokenFromMetadata(context.Background())
	st, _ := status.FromError(err)
	if st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got %v", st.Code())
	}
}
