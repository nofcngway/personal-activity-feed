package auth_action_service_api

import (
	"testing"

	"github.com/nofcngway/auth-action-service/internal/services/authservice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapErr_Codes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		err  error
		code codes.Code
	}{
		{authservice.ErrInvalidArgument, codes.InvalidArgument},
		{authservice.ErrUserAlreadyExists, codes.AlreadyExists},
		{authservice.ErrInvalidCredentials, codes.Unauthenticated},
		{authservice.ErrUnauthorized, codes.Unauthenticated},
		{assertErr{}, codes.Internal},
	}

	for _, tc := range cases {
		st, _ := status.FromError(mapErr(tc.err))
		if st.Code() != tc.code {
			t.Fatalf("for %v expected %v, got %v", tc.err, tc.code, st.Code())
		}
	}
}

type assertErr struct{}

func (assertErr) Error() string { return "x" }


