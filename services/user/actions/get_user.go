package actions

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/planning-pocker/pkg/grpc"
	log "github.com/sweetloveinyourheart/planning-pocker/pkg/logger"
	proto "github.com/sweetloveinyourheart/planning-pocker/proto/code/userserver/go"
)

func (a *actions) GetUser(ctx context.Context, request *connect.Request[proto.GetUserRequest]) (response *connect.Response[proto.GetUserResponse], err error) {
	userID, err := uuid.FromString(request.Msg.GetUserId())
	if err != nil {
		return nil, grpc.InvalidArgumentErrorWithField(grpc.FieldViolation("user_id", fmt.Errorf("user_id is invalid")))
	}

	user, ok, err := a.userRepo.GetUserByID(ctx, userID)
	if !ok {
		log.Global().WarnContext(ctx, "no user was found", zap.Any("userID", userID))
		return nil, err
	}

	if err != nil {
		log.Global().ErrorContext(ctx, "failed to get user data", zap.Error(err))
		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&proto.GetUserResponse{
		User: &proto.User{
			UserId:   user.UserID.String(),
			Username: user.Username,
			FullName: user.FullName,
			Status:   int32(user.Status),
		},
	}), nil
}
