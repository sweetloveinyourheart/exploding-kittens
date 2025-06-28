package actions

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go"
)

func (a *actions) GetUser(ctx context.Context, request *connect.Request[proto.GetUserRequest]) (response *connect.Response[proto.GetUserResponse], err error) {
	opName := "userserver.GetUser()"
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
	}
	ctx, span := a.tracer.Start(ctx, opName, opts...)
	defer span.End()

	userID, err := uuid.FromString(request.Msg.GetUserId())
	if err != nil {
		return nil, grpc.InvalidArgumentErrorWithField(grpc.FieldViolation("user_id", fmt.Errorf("user_id is invalid")))
	}

	user, ok, err := a.userRepo.GetUserByID(ctx, userID)
	if !ok {
		log.Global().WarnContext(ctx, "no user was found", zap.Any("userID", userID))
		return nil, grpc.NotFoundError(errors.New("no user was found"))
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
