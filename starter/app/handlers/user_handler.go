package handlers

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tgo-framework/tgo/pkg/transport/middleware"
	userv1 "github.com/tgo-framework/tgo/starter/proto/v1/userv1"
)

// UserHandler implements the ConnectRPC UserServiceHandler.
type UserHandler struct{}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetProfile handles the GetProfile RPC request.
func (h *UserHandler) GetProfile(
	ctx context.Context,
	req *connect.Request[userv1.GetProfileRequest],
) (*connect.Response[userv1.GetProfileResponse], error) {
	userID := req.Msg.GetUserId()
	if userID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("user_id cannot be empty"))
	}

	tenant := middleware.GetTenant(ctx)
	if tenant == "" {
		tenant = "public"
	}

	resp := &userv1.GetProfileResponse{
		UserId:     userID,
		Name:       fmt.Sprintf("User %s", userID),
		Email:      fmt.Sprintf("user_%s@%s.example.com", userID, tenant),
		Role:       "member",
		TenantSlug: tenant,
	}

	return connect.NewResponse(resp), nil
}
