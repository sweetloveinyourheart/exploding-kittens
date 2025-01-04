package actions_test

import (
	goTesting "testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/planning-pocker/pkg/config"
	"github.com/sweetloveinyourheart/planning-pocker/pkg/testing"
	"github.com/sweetloveinyourheart/planning-pocker/services/user/repos"
	userserver_mock "github.com/sweetloveinyourheart/planning-pocker/services/user/repos/mock"
)

type ActionsSuite struct {
	*testing.Suite
	mockUserRepository           *userserver_mock.MockUserRepository
	mockUserCredentialRepository *userserver_mock.MockUserCredentialRepository
	mockUserSessionRepository    *userserver_mock.MockUserSessionRepository
}

func (as *ActionsSuite) SetupTest() {
	as.mockUserRepository = new(userserver_mock.MockUserRepository)
	as.mockUserCredentialRepository = new(userserver_mock.MockUserCredentialRepository)
	as.mockUserSessionRepository = new(userserver_mock.MockUserSessionRepository)
}

func TestActionsSuite(t *goTesting.T) {
	as := &ActionsSuite{
		Suite:                        testing.MakeSuite(t),
		mockUserRepository:           new(userserver_mock.MockUserRepository),
		mockUserCredentialRepository: new(userserver_mock.MockUserCredentialRepository),
		mockUserSessionRepository:    new(userserver_mock.MockUserSessionRepository),
	}

	suite.Run(t, as)
}

func (as *ActionsSuite) setupEnvironment() {
	do.Override[repos.IUserRepository](nil, func(i *do.Injector) (repos.IUserRepository, error) {
		return as.mockUserRepository, nil
	})

	do.Override[repos.IUserCredentialRepository](nil, func(i *do.Injector) (repos.IUserCredentialRepository, error) {
		return as.mockUserCredentialRepository, nil
	})

	do.Override[repos.IUserSessionRepository](nil, func(i *do.Injector) (repos.IUserSessionRepository, error) {
		return as.mockUserSessionRepository, nil
	})

	config.Instance().Set("userserver.secrets.token_signing_key", "testkey")
}
