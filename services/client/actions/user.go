package actions

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"

	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	userProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
)

func (a *actions) GetUserProfile(ctx context.Context, request *connect.Request[emptypb.Empty]) (response *connect.Response[proto.UserProfileResponse], err error) {
	playerID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		// This should never happen as this endpoint should be authenticated
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	getUserRequest := userProto.GetUserRequest{UserId: playerID.String()}
	profile, err := a.userServerClient.GetUser(ctx, connect.NewRequest(&getUserRequest))
	if err != nil {
		return nil, grpc.NotFoundError(err)
	}

	return connect.NewResponse(&proto.UserProfileResponse{
		User: &proto.User{
			UserId:   profile.Msg.GetUser().GetUserId(),
			Username: profile.Msg.GetUser().GetUsername(),
			FullName: profile.Msg.GetUser().GetFullName(),
			Status:   profile.Msg.GetUser().GetStatus(),
		},
	}), nil
}

func (a *actions) GetPlayersProfile(ctx context.Context, request *connect.Request[proto.PlayersProfileRequest]) (response *connect.Response[proto.PlayersProfileResponse], err error) {
	_, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		// This should never happen as this endpoint should be authenticated
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	profiles := make([]*proto.User, 0)
	for _, userID := range request.Msg.GetUserIds() {
		getUserRequest := userProto.GetUserRequest{UserId: userID}
		profile, err := a.userServerClient.GetUser(ctx, connect.NewRequest(&getUserRequest))
		if err != nil {
			return nil, grpc.NotFoundError(err)
		}

		user := &proto.User{
			UserId:   profile.Msg.GetUser().GetUserId(),
			Username: profile.Msg.GetUser().GetUsername(),
			FullName: profile.Msg.GetUser().GetFullName(),
			Status:   profile.Msg.GetUser().GetStatus(),
		}
		profiles = append(profiles, user)
	}

	return connect.NewResponse(&proto.PlayersProfileResponse{
		Users: profiles,
	}), nil
}
