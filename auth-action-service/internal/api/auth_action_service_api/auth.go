package auth_action_service_api

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func tokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing authorization")
	}
	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization")
	}
	h := strings.TrimSpace(values[0])
	if h == "" {
		return "", status.Error(codes.Unauthenticated, "missing authorization")
	}

	// Поддерживаем и "Bearer <token>", и просто "<token>".
	lh := strings.ToLower(h)
	const prefix = "bearer "
	if strings.HasPrefix(lh, prefix) {
		h = strings.TrimSpace(h[len(prefix):])
	}
	if h == "" {
		return "", status.Error(codes.Unauthenticated, "invalid authorization header")
	}
	return h, nil
}
