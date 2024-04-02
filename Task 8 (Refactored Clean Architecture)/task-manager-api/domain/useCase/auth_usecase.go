package usecase

import "task-manager-api/model"

type AuthUseCase interface {
	RegisterUser( currUser *model.AuthenticatedUser, userCreate *model.UserCreate, param any) (*model.UserInfo, string, error)
	Login( currUser *model.AuthenticatedUser, userLogin *model.UserLogin, param any) (*model.Token, string, error)
	RegisterAdmin( currUser *model.AuthenticatedUser, userCreate *model.UserCreate, param any) (*model.UserInfo, string, error)
}