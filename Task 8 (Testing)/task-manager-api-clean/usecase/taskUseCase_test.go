package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"task-manager-api-clean/domain"
	"task-manager-api-clean/domain/mocks"
    "task-manager-api-clean/usecase"


	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
    testId         = "testId"
    nonExistentId  = "nonExistentId"
)

type TaskUseCaseTestSuite struct {
    suite.Suite
    repo    *mocks.TaskRepository
    useCase domain.TaskUseCase

}

func (suite *TaskUseCaseTestSuite) SetupTest() {
	suite.repo = new(mocks.TaskRepository)
    suite.useCase = usecase.NewTaskUseCase(suite.repo)

}

func (suite *TaskUseCaseTestSuite) TestCreate_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    taskInput := &domain.TaskInput{Title: "Test Task", Description: "Test Description", Status: "Pending", DueDate: time.Now()}
    task := &domain.Task{Title: taskInput.Title, Description: taskInput.Description, Status: taskInput.Status, DueDate: taskInput.DueDate}

    suite.repo.On("Create", ctx, mock.Anything).Return(task, nil)

    createdTask, err := suite.useCase.Create(ctx, taskInput)
    suite.NoError(err)
    suite.Equal(taskInput.Title, createdTask.Title)
}

func (suite *TaskUseCaseTestSuite) TestCreate_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    taskInput := &domain.TaskInput{Title: "Test Task", Description: "Test Description", Status: "Pending", DueDate: time.Now()}

    suite.repo.On("Create", ctx, mock.Anything).Return(nil, errors.New("create error"))

    _, err := suite.useCase.Create(ctx, taskInput)
    suite.Error(err)
}

func (suite *TaskUseCaseTestSuite) TestUpdate_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    taskInput := &domain.TaskInput{Title: "Updated Task", Description: "Updated Description", Status: "Completed", DueDate: time.Now()}
    task := &domain.Task{Id: testId, Title: "Test Task", Description: "Test Description", Status: "Pending", DueDate: time.Now()}
    updatedTask := &domain.Task{Id: task.Id, Title: taskInput.Title, Description: taskInput.Description, Status: taskInput.Status, DueDate: taskInput.DueDate}

    suite.repo.On("GetById", ctx, testId).Return(task, nil)
    suite.repo.On("Update", ctx, testId, mock.Anything).Return(updatedTask, nil)

    result, err := suite.useCase.Update(ctx, testId, taskInput)
    suite.NoError(err)
    suite.Equal(taskInput.Title, result.Title)
}

func (suite *TaskUseCaseTestSuite) TestUpdate_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    taskInput := &domain.TaskInput{Title: "Updated Task", Description: "Updated Description", Status: "Completed", DueDate: time.Now()}

    suite.repo.On("GetById", ctx, nonExistentId).Return(nil, errors.New("not found"))

    _, err := suite.useCase.Update(ctx, nonExistentId, taskInput)
    suite.Error(err)
}

func (suite *TaskUseCaseTestSuite) TestDelete_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    suite.repo.On("Delete", ctx, testId).Return(nil)

    err := suite.useCase.Delete(ctx, testId)
    suite.NoError(err)
}

func (suite *TaskUseCaseTestSuite) TestDelete_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    suite.repo.On("Delete", ctx, nonExistentId).Return(errors.New("delete error"))

    err := suite.useCase.Delete(ctx, nonExistentId)
    suite.Error(err)
}

func (suite *TaskUseCaseTestSuite) TestGetAll_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    tasks := []*domain.Task{
        {Title: "Test Task 1", Description: "Test Description 1", Status: "Pending", DueDate: time.Now()},
        {Title: "Test Task 2", Description: "Test Description 2", Status: "Completed", DueDate: time.Now()},
    }

    suite.repo.On("GetAll", ctx).Return(&tasks, nil)

    result, err := suite.useCase.GetAll(ctx)
    suite.NoError(err)
    suite.Len(*result, 2)
}

func (suite *TaskUseCaseTestSuite) TestGetAll_Empty() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var tasks []*domain.Task

    suite.repo.On("GetAll", ctx).Return(&tasks, nil)

    result, err := suite.useCase.GetAll(ctx)
    suite.NoError(err)
    suite.Len(*result, 0)
}

func (suite *TaskUseCaseTestSuite) TestGetById_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    task := &domain.Task{Id: testId, Title: "Test Task", Description: "Test Description", Status: "Pending", DueDate: time.Now()}

    suite.repo.On("GetById", ctx, testId).Return(task, nil)

    result, err := suite.useCase.GetById(ctx, testId)
    suite.NoError(err)
    suite.Equal(task.Title, result.Title)
}

func (suite *TaskUseCaseTestSuite) TestGetById_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    suite.repo.On("GetById", ctx, nonExistentId).Return(nil, errors.New("not found"))

    _, err := suite.useCase.GetById(ctx, nonExistentId)
    suite.Error(err)
}

func TestTaskUseCaseTestSuite(t *testing.T) {
    suite.Run(t, new(TaskUseCaseTestSuite))
}