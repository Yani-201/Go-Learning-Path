package usecase

import (
	"task-manager-api/model"
)
type UserUseCase interface {
    GetUsers(currUser *model.AuthenticatedUser, dto any, param *model.SearchParam) (*[]*model.UserInfo, string, error)
    GetUserByID(currUser *model.AuthenticatedUser, dto any, param *model.IdParam) (*model.UserInfo, string, error)
	DeleteUser(currUser *model.AuthenticatedUser, dto any, param *model.IdParam) (*model.UserInfo, string, error)
	UpdateUser(currUser *model.AuthenticatedUser, updated *model.UserUpdate, param any) (*model.UserInfo, string, error)
	UpdatePassword(currUser *model.AuthenticatedUser, updated *model.UserUpdatePassword,param any) (*model.UserInfo, string, error)

}