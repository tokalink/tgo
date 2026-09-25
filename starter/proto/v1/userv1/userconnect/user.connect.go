package userconnect

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	userv1 "github.com/tokalink/tgo/starter/proto/v1/userv1"
)

const (
	// UserServiceName is the fully-qualified name of the UserService service.
	UserServiceName = "user.v1.UserService"
	// UserServiceGetProfileProcedure is the fully-qualified name of the UserService's GetProfile RPC.
	UserServiceGetProfileProcedure = "/user.v1.UserService/GetProfile"
)

// UserServiceHandler is an implementation of the user.v1.UserService service.
type UserServiceHandler interface {
	GetProfile(context.Context, *connect.Request[userv1.GetProfileRequest]) (*connect.Response[userv1.GetProfileResponse], error)
}

// NewUserServiceHandler builds an HTTP handler from the service implementation.
// It returns the path on which to mount the handler and the handler itself.
func NewUserServiceHandler(svc UserServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	userServiceGetProfileHandler := connect.NewUnaryHandler(
		UserServiceGetProfileProcedure,
		svc.GetProfile,
		opts...,
	)

	mux := http.NewServeMux()
	mux.Handle(UserServiceGetProfileProcedure, userServiceGetProfileHandler)
	return "/user.v1.UserService/", mux
}

// UserServiceClient is a client for the user.v1.UserService service.
type UserServiceClient interface {
	GetProfile(context.Context, *connect.Request[userv1.GetProfileRequest]) (*connect.Response[userv1.GetProfileResponse], error)
}

type userServiceClient struct {
	getProfile *connect.Client[userv1.GetProfileRequest, userv1.GetProfileResponse]
}

// NewUserServiceClient constructs a client for the user.v1.UserService service.
func NewUserServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) UserServiceClient {
	return &userServiceClient{
		getProfile: connect.NewClient[userv1.GetProfileRequest, userv1.GetProfileResponse](
			httpClient,
			baseURL+UserServiceGetProfileProcedure,
			opts...,
		),
	}
}

func (c *userServiceClient) GetProfile(ctx context.Context, req *connect.Request[userv1.GetProfileRequest]) (*connect.Response[userv1.GetProfileResponse], error) {
	if c.getProfile == nil {
		return nil, errors.New("client not initialized")
	}
	return c.getProfile.CallUnary(ctx, req)
}
