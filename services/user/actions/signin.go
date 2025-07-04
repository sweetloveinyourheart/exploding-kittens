package actions

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/otel/trace"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/user/models"
	"github.com/sweetloveinyourheart/exploding-kittens/services/user/utils"
)

func (a *actions) SignIn(ctx context.Context, request *connect.Request[proto.SignInRequest]) (response *connect.Response[proto.SignInResponse], err error) {
	opName := "userserver.SignIn()"
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
	}
	ctx, span := a.tracer.Start(ctx, opName, opts...)
	defer span.End()

	userID, err := uuid.FromString(request.Msg.GetUserId())
	if err != nil {
		return nil, grpc.InvalidArgumentErrorWithField(grpc.FieldViolation("user_id", err))
	}

	// Lookup user
	user, found, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, grpc.UnauthenticatedError(err)
	}
	if !found {
		err = errors.New("unauthorized, user not found")
		return nil, grpc.NotFoundError(err)
	}

	// Build session
	sessionHash, err := utils.SessionHash(userID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	now := time.Now()
	expiration := now.Add(time.Hour * 24)

	session := models.UserSession{
		SessionID:         sessionID,
		Token:             sessionHash,
		UserID:            userID,
		SessionStart:      now,
		LastUpdated:       now,
		SessionExpiration: &expiration,
	}

	err = a.userSessionRepo.CreateSession(ctx, &session)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	// Sign token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID.String(),
		"token":   session.Token,
	})

	signingKey := config.Instance().GetString("userserver.secrets.token_signing_key")
	if len(signingKey) == 0 {
		err = errors.New("invalid signing token, misconfigured instance")
		return nil, grpc.InternalError(err)
	}
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&proto.SignInResponse{
		User: &proto.User{
			UserId:    user.UserID.String(),
			Username:  user.Username,
			FullName:  user.FullName,
			Status:    int32(user.Status),
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		},
		Token: tokenString,
	}), nil
}
