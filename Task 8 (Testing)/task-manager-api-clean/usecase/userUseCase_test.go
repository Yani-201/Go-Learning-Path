package usecase_test

import (
	"context"
	"errors"
	"testing"

	"task-manager-api-clean/config"
	"task-manager-api-clean/domain"
	"task-manager-api-clean/domain/mocks"
    "task-manager-api-clean/usecase"
	"task-manager-api-clean/utils"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUseCaseTestSuite struct {
    suite.Suite
    repo    *mocks.UserRepository
    useCase domain.UserUseCase
    env     *config.Environment
}

func (suite *UserUseCaseTestSuite) SetupTest() {
    suite.repo = new(mocks.UserRepository)
	suite.env = &config.Environment{JwtSecret: "secret"}
    suite.useCase = usecase.NewUserUseCase(suite.repo, suite.env)

}

func (suite *UserUseCaseTestSuite) TestRegisterUser_Success() {
    payload := &domain.UserCreate{Username: "test", Password: "test", Email: "test@example.com"}
    suite.repo.On("GetByUsername", mock.Anything, "test").Return(nil, errors.New("user not found"))
    suite.repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(&domain.User{UserID: "1", Username: "test", Email: "test@example.com"}, nil)

    userInfo, err := suite.useCase.RegisterUser(context.Background(), payload)
    suite.NoError(err)

	expectedUserInfo := &domain.UserInfo{
		Username: "test",
		Email: "test@example.com",
		UserId: "1",
	}
	
	suite.Equal(expectedUserInfo, userInfo)
    suite.Equal(expectedUserInfo, userInfo)
    suite.repo.AssertCalled(suite.T(), "GetByUsername", mock.Anything, "test")
    suite.repo.AssertCalled(suite.T(), "Create", mock.Anything, mock.AnythingOfType("*domain.User"))
}

func (suite *UserUseCaseTestSuite) TestRegisterUser_EmptyUsername() {
    payload := &domain.UserCreate{Username: "", Password: "test", Email: "test@example.com"}

    _, err := suite.useCase.RegisterUser(context.Background(), payload)
    suite.EqualError(err, "invalid Payload")
}

func (suite *UserUseCaseTestSuite) TestRegisterUser_UserExists() {
    payload := &domain.UserCreate{Username: "test", Password: "test", Email: "test@example.com"}
    suite.repo.On("GetByUsername", mock.Anything, "test").Return(&domain.User{UserID: "1", Username: "test", Email: "test@example.com"}, nil)

    _, err := suite.useCase.RegisterUser(context.Background(), payload)
    suite.EqualError(err, "user already exists")
}

func (suite *UserUseCaseTestSuite) TestLogin_Success() {
    payload := &domain.UserLogin{Username: "test", Password: "hashedPassword"}
    hashedPassword, err := utils.EncryptPassword("hashedPassword")
    suite.NoError(err)
    
    suite.repo.On("GetByUsername", mock.Anything, "test").Return(&domain.User{UserID: "1", Username: "test", Password: hashedPassword}, nil)
    suite.repo.On("GetByUsername", mock.Anything, "test").Return(&domain.User{Username: "test", Password: hashedPassword}, nil)

    jwtToken, err := suite.useCase.Login(context.Background(), payload)
    suite.NoError(err)
    suite.NotEmpty(jwtToken)
    suite.repo.AssertCalled(suite.T(), "GetByUsername", mock.Anything, "test")
}

func (suite *UserUseCaseTestSuite) TestLogin_InvalidCredentials() {
    payload := &domain.UserLogin{Username: "test", Password: "wrongpassword"}
    hashedPassword, err := utils.EncryptPassword("hashedPassword")
    suite.NoError(err)
    
    suite.repo.On("GetByUsername", mock.Anything, "test").Return(&domain.User{UserID: "1", Username: "test", Password: hashedPassword}, nil)

    _, err = suite.useCase.Login(context.Background(), payload)
    suite.EqualError(err, "invalid username or password")
}

func (suite *UserUseCaseTestSuite) TestPromote_Success() {
    username := "testUser"
    suite.repo.On("GetByUsername", mock.Anything, username).Return(&domain.User{UserID: "1", Username: username, Role: "user"}, nil)
    suite.repo.On("UpdateRole", mock.Anything, username, "admin").Return(&domain.User{UserID: "1", Username: username, Role: "admin"}, nil)

    _, err := suite.useCase.Promote(context.Background(), username)
    suite.NoError(err)
    
    suite.repo.AssertCalled(suite.T(), "GetByUsername", mock.Anything, username)
    suite.repo.AssertCalled(suite.T(), "UpdateRole", mock.Anything, username, "admin")
}

func (suite *UserUseCaseTestSuite) TestPromote_UserNotFound() {
    username := "nonExistentUser"
    suite.repo.On("GetByUsername", mock.Anything, username).Return(nil, errors.New("user not found"))

    _, err := suite.useCase.Promote(context.Background(), username)
    suite.EqualError(err, "user not found")
}

func TestUserUseCaseTestSuite(t *testing.T) {
    suite.Run(t, new(UserUseCaseTestSuite))
}