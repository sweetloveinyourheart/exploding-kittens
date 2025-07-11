package actions

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/user/models"
)

func (a *actions) CreateNewUser(ctx context.Context, request *connect.Request[proto.CreateUserRequest]) (response *connect.Response[proto.CreateUserResponse], err error) {
	opName := "userserver.CreateNewUser()"
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
	}
	ctx, span := a.tracer.Start(ctx, opName, opts...)
	defer span.End()

	newUser := models.User{
		UserID:    uuid.Must(uuid.NewV7()),
		Username:  request.Msg.GetUsername(),
		FullName:  request.Msg.GetFullName(),
		Status:    models.USER_STATUS_ENABLED,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = newUser.Validate()
	if err != nil {
		return nil, grpc.InvalidArgumentError(err)
	}

	_, found, err := a.userRepo.GetUserByUsername(ctx, request.Msg.GetUsername())
	if err != nil {
		return nil, grpc.InternalError(err)
	}
	if found {
		err = errors.New("user already exists !")
		log.Global().ErrorContext(ctx, err.Error(), zap.Error(err))
		return nil, grpc.AlreadyExistsError(err)
	}

	// Create new user
	err = a.userRepo.CreateUser(ctx, &newUser)
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to create new user", zap.Error(err))
		return nil, grpc.InternalError(err)
	}

	// Create new credential
	var meta []byte
	if request.Msg.GetMeta() != "" {
		meta = []byte(request.Msg.GetMeta())
	}

	newUserCredential := models.UserCredential{
		UserID:       newUser.UserID,
		AuthProvider: request.Msg.GetAuthProvider().String(),
		Meta:         meta,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = newUserCredential.Validate()
	if err != nil {
		return nil, grpc.InvalidArgumentError(err)
	}

	err = a.userCredentialRepo.CreateCredential(ctx, &newUserCredential)
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to create new user credential", zap.Error(err))
		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&proto.CreateUserResponse{
		User: &proto.User{
			UserId:    newUser.UserID.String(),
			Username:  newUser.Username,
			FullName:  newUser.FullName,
			Status:    int32(newUser.Status),
			CreatedAt: newUser.CreatedAt.Unix(),
			UpdatedAt: newUser.UpdatedAt.Unix(),
		},
	}), nil
}
