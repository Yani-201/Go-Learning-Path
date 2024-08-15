package usecase

import (
	"context"
	"task-manager-api-clean/domain"
	
)


type TaskUseCase struct {
	TaskRepository domain.TaskRepository

}

func NewTaskUseCase(tr domain.TaskRepository) domain.TaskUseCase {
	return &TaskUseCase{
		TaskRepository: tr,
	}
}
func (tu *TaskUseCase) Create(c context.Context, payload *domain.TaskInput) (*domain.Task, error) {
	task := &domain.Task{
		Title:       payload.Title,
		Description: payload.Description,
		Status:      payload.Status,
		DueDate:     payload.DueDate,
	}
	
	return tu.TaskRepository.Create(c, task)
}

func (tu *TaskUseCase) Update(c context.Context, taskId string, payload *domain.TaskInput) (*domain.Task, error) {
	task, err := tu.TaskRepository.GetById(c, taskId)
	if err != nil {
		return nil, err
	}

	if payload.Title != "" {
		task.Title = payload.Title
	}
	if payload.Description != "" {
		task.Description = payload.Description
	}
	if !payload.DueDate.IsZero() {
		task.DueDate = payload.DueDate
	}
	if task.Status != payload.Status {
		task.Status = payload.Status
	}

	return tu.TaskRepository.Update(c, taskId, task)
}

func (tu *TaskUseCase) Delete(c context.Context, taskId string) error {
	return tu.TaskRepository.Delete(c, taskId)
}

func (tu *TaskUseCase) GetAll(c context.Context) (*[]*domain.Task, error) {
	return tu.TaskRepository.GetAll(c)
}

func (tu *TaskUseCase) GetById(c context.Context, taskId string) (*domain.Task, error) {
	return tu.TaskRepository.GetById(c, taskId)
}
