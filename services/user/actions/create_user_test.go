package actions_test

import (
	"context"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"

	proto "github.com/sweetloveinyourheart/planning-pocker/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/planning-pocker/services/user/actions"
	"github.com/sweetloveinyourheart/planning-pocker/services/user/models"
)

func (as *ActionsSuite) Test_CreateNewUser_NoUsername() {
	as.setupEnvironment()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	actions := actions.NewActions(ctx, "test")
	resp, err := actions.CreateNewUser(
		ctx,
		connect.NewRequest(
			&proto.CreateUserRequest{
				FullName: "John Due",
			},
		),
	)

	as.ErrorContains(err, "Username: blank")
	as.Nil(resp)
}

func (as *ActionsSuite) Test_CreateNewUser_UserIsExisting() {
	as.setupEnvironment()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	as.mockUserRepository.On("GetUserByUsername", mock.Anything, mock.Anything).Return(models.User{Username: "John"}, true, nil)

	newUser := &proto.CreateUserRequest{
		Username:     "John",
		FullName:     "John Due",
		AuthProvider: proto.CreateUserRequest_GUEST,
		Meta:         nil,
	}

	actions := actions.NewActions(ctx, "test")
	resp, err := actions.CreateNewUser(
		ctx,
		connect.NewRequest(newUser),
	)

	as.Nil(resp)
	as.Error(err)

	as.mockUserRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) Test_CreateNewUser_Success() {
	as.setupEnvironment()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	as.mockUserRepository.On("GetUserByUsername", mock.Anything, mock.Anything).Return(models.User{}, false, nil)
	as.mockUserRepository.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
	as.mockUserCredentialRepository.On("CreateCredential", mock.Anything, mock.Anything).Return(nil)

	newUser := &proto.CreateUserRequest{
		Username:     "John",
		FullName:     "John Due",
		AuthProvider: proto.CreateUserRequest_GUEST,
		Meta:         nil,
	}

	actions := actions.NewActions(ctx, "test")
	resp, err := actions.CreateNewUser(
		ctx,
		connect.NewRequest(newUser),
	)

	as.Nil(err)
	as.Equal(newUser.Username, resp.Msg.User.GetUsername())
	as.Equal(newUser.FullName, resp.Msg.User.GetFullName())
	as.Equal(int32(models.USER_STATUS_ENABLED), resp.Msg.User.GetStatus())

	as.mockUserRepository.AssertExpectations(as.T())
	as.mockUserCredentialRepository.AssertExpectations(as.T())
}
