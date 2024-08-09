package usecase

import (
	"task-manager-api/domain"
	"task-manager-api/model"
)

type TaskUseCase interface {

	CreateTask(currUser *model.AuthenticatedUser,dto *model.TaskCreate, param any) (*model.TaskInfo, string, error)
	GetTasks(currUser *model.AuthenticatedUser, dto any, param *model.SearchParam) (*[]*model.TaskInfo, string, error)
	GetTasksByID(currUser *model.AuthenticatedUser, dto any, param *model.IdParam) (*model.TaskInfo, string, error)
	GetTasksByUserID(currUser *model.AuthenticatedUser, dto any, param *model.IdParam) (*model.TaskInfo, string, error)
	UpdateTask(currUser *model.AuthenticatedUser, dto *model.TaskUpdate, param *model.IdParam) (*model.TaskInfo, string, error)
	DeleteTask(currUser *model.AuthenticatedUser, dto any, param *model.IdParam) (*domain.Task, string, error)
}