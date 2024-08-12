package usecase

import (
	"context"
	"errors"

	"task-manager-api-clean/utils"
	"task-manager-api-clean/domain"
	"task-manager-api-clean/config"

)


type UserUseCase struct {
	environment *config.Environment
	UserRepository domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository, env *config.Environment) domain.UserUseCase {
	return &UserUseCase{
		UserRepository: userRepo,
		environment: env,
	}
}

func (uc *UserUseCase) RegisterUser(c context.Context, payload *domain.UserCreate) (*domain.UserInfo, error) {
	if payload.Username == "" || payload.Password == "" {
		return nil, errors.New("invalid Payload")
	}
	
	user := &domain.User{
		Username: payload.Username,
		Password: payload.Password,
		Email:    payload.Email,
	}
	//check if user exists
	_, err := uc.UserRepository.GetByUsername(c, user.Username)
	if err == nil {
		return nil, errors.New("user already exists")	
	}
	

	createdUser, err := uc.UserRepository.Create(c, user)
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		UserId:   createdUser.UserID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
	}, nil
}

func (uc *UserUseCase) Login(c context.Context, payload *domain.UserLogin) (string, error) {
	user, err := uc.UserRepository.GetByUsername(c, payload.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Compare passwords
	err = utils.ComparePasswords(user.Password, payload.Password)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT token
	jwtToken, err := utils.TokenGenerate(&domain.AuthenticatedUser{
		UserID:   user.UserID,
		Username: user.Username,
		Role:     user.Role,
	}, uc.environment.JwtSecret)

	return jwtToken, err
}

func (uc *UserUseCase) Promote(c context.Context, username string) (*domain.UserInfo, error) {
	user, err := uc.UserRepository.GetByUsername(c, username)
    if err != nil {
        return nil, err
    }

    if user.Role == "admin" {
        return nil, errors.New("user is already an admin")
    }

	user, err = uc.UserRepository.UpdateRole(c, username, "admin")
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		UserId:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
