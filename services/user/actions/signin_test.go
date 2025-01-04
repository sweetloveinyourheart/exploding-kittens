package actions_test

import (
	"context"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	proto "github.com/sweetloveinyourheart/planning-pocker/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/planning-pocker/services/user/actions"
	"github.com/sweetloveinyourheart/planning-pocker/services/user/models"
)

func (as *ActionsSuite) Test_SignIn_NoUserWasFound() {
	as.setupEnvironment()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	as.mockUserRepository.On("GetUserByID", mock.Anything, mock.Anything).Return(models.User{}, false, nil)

	userID := uuid.Must(uuid.NewV7())

	actions := actions.NewActions(ctx, "test")
	resp, err := actions.SignIn(
		ctx,
		connect.NewRequest(
			&proto.SignInRequest{
				UserId: userID.String(),
			},
		),
	)

	as.Nil(resp)
	as.ErrorContains(err, "unauthorized, user not found")

	as.mockUserRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) Test_SignIn_Success() {
	as.setupEnvironment()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	as.mockUserRepository.On("GetUserByID", mock.Anything, mock.Anything).Return(
		models.User{
			UserID:   uuid.Must(uuid.NewV7()),
			Username: "john",
			FullName: "John Due",
			Status:   1,
		},
		true,
		nil,
	)
	as.mockUserSessionRepository.On("CreateSession", mock.Anything, mock.Anything).Return(nil)

	userID := uuid.Must(uuid.NewV7())

	actions := actions.NewActions(ctx, "test")
	_, err := actions.SignIn(
		ctx,
		connect.NewRequest(
			&proto.SignInRequest{
				UserId: userID.String(),
			},
		),
	)

	as.Nil(err)

	as.mockUserRepository.AssertExpectations(as.T())
	as.mockUserSessionRepository.AssertExpectations(as.T())
}
