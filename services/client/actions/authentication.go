package actions

import (
	"context"

	"connectrpc.com/connect"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	userProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go"
)

func (a *actions) CreateNewGuestUser(ctx context.Context, request *connect.Request[proto.CreateNewGuestUserRequest]) (response *connect.Response[proto.CreateNewGuestUserResponse], err error) {
	newGuestUser := userProto.CreateUserRequest{
		Username:     request.Msg.GetUsername(),
		FullName:     request.Msg.GetFullName(),
		AuthProvider: userProto.CreateUserRequest_GUEST,
		Meta:         nil,
	}

	resp, err := a.userServerClient.CreateNewUser(ctx, connect.NewRequest(&newGuestUser))
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&proto.CreateNewGuestUserResponse{
		User: &proto.User{
			UserId:   resp.Msg.GetUser().GetUserId(),
			Username: resp.Msg.GetUser().GetUsername(),
			FullName: resp.Msg.GetUser().GetFullName(),
			Status:   resp.Msg.GetUser().GetStatus(),
		},
	}), nil
}

func (a *actions) GuestLogin(ctx context.Context, request *connect.Request[proto.GuestLoginRequest]) (response *connect.Response[proto.GuestLoginResponse], err error) {
	credentials := userProto.SignInRequest{
		UserId: request.Msg.GetUserId(),
	}

	resp, err := a.userServerClient.SignIn(ctx, connect.NewRequest(&credentials))
	if err != nil {
		return nil, grpc.UnauthenticatedError(err)
	}

	return connect.NewResponse(&proto.GuestLoginResponse{
		User: &proto.User{
			UserId:   resp.Msg.GetUser().GetUserId(),
			Username: resp.Msg.GetUser().GetUsername(),
			FullName: resp.Msg.GetUser().GetFullName(),
			Status:   resp.Msg.GetUser().GetStatus(),
		},
		Token: resp.Msg.GetToken(),
	}), nil
}
